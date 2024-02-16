package core

import (
	"context"
	"database/sql"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	querierPkg "github.com/HomayoonAlimohammadi/order-delay/internal/database/querier"
	orderpb "github.com/HomayoonAlimohammadi/order-delay/proto"
)

type TripsHandler interface {
	IsOngoing(trip querierPkg.Trip) bool
}

type DeliveryTimeEstimator interface {
	Estimate(ctx context.Context, order querierPkg.Order) (time.Duration, error)
}

type Logger interface {
	Errorf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Infof(format string, args ...interface{})
}

type PgxPool interface {
	querierPkg.DBTX
	Close()
}

type DelayCheckQueue interface {
	Put(ctx context.Context, ordr querierPkg.Order) error
}

type orderServiceImpl struct {
	orderpb.UnimplementedOrderDelayServer

	dbpool  PgxPool
	logger  Logger
	querier querierPkg.Querier

	tripsHandler          TripsHandler
	deliveryTimeEstimator DeliveryTimeEstimator
	delayCheckQueue       DelayCheckQueue
	nowFunc               func() time.Time
}

func New(
	dbpool PgxPool,
	logger Logger,
	querier querierPkg.Querier,
	tripsHandler TripsHandler,
	deliveryTimeEstimator DeliveryTimeEstimator,
	delayCheckQueue DelayCheckQueue,
	nowFunc func() time.Time,
) *orderServiceImpl {
	return &orderServiceImpl{
		dbpool:                dbpool,
		logger:                logger,
		querier:               querier,
		tripsHandler:          tripsHandler,
		deliveryTimeEstimator: deliveryTimeEstimator,
		delayCheckQueue:       delayCheckQueue,
		nowFunc:               nowFunc,
	}
}

func (s *orderServiceImpl) ReportOrderDelay(ctx context.Context, req *orderpb.ReportOrderDelayRequest) (*orderpb.ReportOrderDelayResponse, error) {
	// create delay report
	err := s.querier.CreateDelayReport(ctx, req.OrderId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create delay report for order: %s", err)
	}

	ordr, err := s.querier.GetOrderByID(ctx, req.OrderId)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "failed to get order by id: %s", err)
	}

	// try to get trip for this order
	// if not found, put order in delay queue and return
	trip, err := s.querier.GetTripByOrderID(ctx, req.OrderId)
	if err != nil {
		// not found
		if err == sql.ErrNoRows {
			err = s.delayCheckQueue.Put(ctx, ordr)
			if err != nil {
				return nil, status.Errorf(codes.Unknown, "failed to put order in delay queue: %s", err)
			}

			return &orderpb.ReportOrderDelayResponse{Result: &orderpb.ReportOrderDelayResponse_Status_{
				Status: orderpb.ReportOrderDelayResponse_REPORTED,
			}}, nil
		}

		return nil, status.Errorf(codes.Unknown, "failed to get trip by order ID: %s", err)
	}

	// if the trip is ongoing, return a new delivery time estimation
	if s.tripsHandler.IsOngoing(trip) {
		newDeliveryTime, err := s.deliveryTimeEstimator.Estimate(ctx, ordr)
		if err != nil {
			return nil, status.Errorf(codes.Unknown, "failed to get new delivery time from estimator: %s", err)
		}

		err = s.querier.UpdateOrderDeliveryTime(ctx,
			querierPkg.UpdateOrderDeliveryTimeParams{
				ID:           req.OrderId,
				DeliveryTime: pgtype.Interval{Microseconds: newDeliveryTime.Microseconds(), Valid: true},
			})
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to update delivery time: %s", err)
		}

		return &orderpb.ReportOrderDelayResponse{Result: &orderpb.ReportOrderDelayResponse_NewDeliveryTimeMinutes{
			NewDeliveryTimeMinutes: int32(newDeliveryTime.Minutes()),
		}}, nil
	}

	// if the trip was not ongoing, file a delay report
	err = s.delayCheckQueue.Put(ctx, ordr)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "failed to put order in delay queue: %s", err)
	}

	return &orderpb.ReportOrderDelayResponse{Result: &orderpb.ReportOrderDelayResponse_Status_{
		Status: orderpb.ReportOrderDelayResponse_REPORTED,
	}}, nil
}

func (s *orderServiceImpl) AssignDelayReportToAgent(ctx context.Context, req *orderpb.AssignDelayReportToAgentRequest) (*emptypb.Empty, error) {
	err := s.querier.CreateDelayCheck(ctx, querierPkg.CreateDelayCheckParams{
		AgentID:  req.AgentId,
		ReportID: req.ReportId,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create delay report check for agent: %s", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *orderServiceImpl) GetVendorDelays(ctx context.Context, req *orderpb.GetVendorDelaysRequest) (*orderpb.GetVendorDelaysResponse, error) {
	oneWeekAgo := s.nowFunc().Add(-7 * 24 * time.Hour)
	reports, err := s.querier.ListDelayReportsForVendor(ctx, querierPkg.ListDelayReportsForVendorParams{
		VendorID:  req.VendorId,
		CreatedAt: pgtype.Timestamptz{Time: oneWeekAgo},
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return &orderpb.GetVendorDelaysResponse{}, nil // no reports
		}

		return nil, status.Errorf(codes.Internal, "failed to get vendor delays: %s", err)
	}

	reportspb := make([]*orderpb.DelayReport, len(reports))
	for i, r := range reports {
		reportspb[i] = &orderpb.DelayReport{Id: r.ID, OrderId: r.OrderID}
	}

	return &orderpb.GetVendorDelaysResponse{DelayReports: reportspb}, nil
}

// for integration testing purposes, e.g. with dockertest
func (s *orderServiceImpl) Close() {
	s.dbpool.Close()
}
