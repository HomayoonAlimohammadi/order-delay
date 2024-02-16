package main

import (
	"context"
	"net"
	"net/url"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/HomayoonAlimohammadi/order-delay/internal/clients"
	"github.com/HomayoonAlimohammadi/order-delay/internal/config"
	"github.com/HomayoonAlimohammadi/order-delay/internal/core"
	querierPkg "github.com/HomayoonAlimohammadi/order-delay/internal/database/querier"
	"github.com/HomayoonAlimohammadi/order-delay/internal/trips"
	orderpb "github.com/HomayoonAlimohammadi/order-delay/proto"
)

const defaultPort = "8888"

func main() {
	cfg := config.Load()
	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}

	pgURL := os.Getenv("POSTGRES_URL")
	if pgURL == "" {
		logger.Fatal("POSTGRES_URL must be set")
	}

	parsedURL, err := url.Parse(pgURL)
	if err != nil {
		logger.WithError(err).Fatal("Failed to parse POSTGRES_URL")
	}

	port := defaultPort
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logger.WithError(err).Fatal("Failed to create listener")
	}

	s := grpc.NewServer()
	reflection.Register(s)

	dbpool, err := pgxpool.New(context.Background(), parsedURL.String())
	if err != nil {
		logger.WithError(err).Fatal("failed to start new pgx pool")
	}
	defer dbpool.Close()

	tripsHandler := getTripsHandler()
	deliveryTimeEstimator := getDeliveryTimeEstimator(cfg.Clients.DeliveryTimeEstimator)
	delayCheckQueue := getDelayCheckQueue(cfg.Clients.DelayCheckQueue)

	srv := core.New(
		dbpool,
		logger,
		querierPkg.New(dbpool),
		tripsHandler,
		deliveryTimeEstimator,
		delayCheckQueue,
		time.Now,
	)

	logger.Info("starting the server...")

	orderpb.RegisterOrderDelayServer(s, srv)
	err = s.Serve(lis)
	if err != nil {
		logger.WithError(err).Fatal("Failed to serve grpc server")
	}
}

func getTripsHandler() core.TripsHandler {
	return trips.NewHandler()
}

func getDeliveryTimeEstimator(c config.DeliveryTimeEstimatorConfig) core.DeliveryTimeEstimator {
	return clients.NewDeliveryTimeEstimator(c.Timeout)
}

func getDelayCheckQueue(c config.DelayCheckQueueConfig) core.DelayCheckQueue {
	return clients.NewDelayCheckQueue(c.Timeout)
}
