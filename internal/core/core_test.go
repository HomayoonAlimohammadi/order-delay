package core_test

import (
	"context"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/HomayoonAlimohammadi/order-delay/internal/core"
	coreMocks "github.com/HomayoonAlimohammadi/order-delay/internal/core/mocks"
	"github.com/HomayoonAlimohammadi/order-delay/internal/database/querier"
	querierMocks "github.com/HomayoonAlimohammadi/order-delay/internal/database/querier/mocks"
	orderpb "github.com/HomayoonAlimohammadi/order-delay/proto"
)

type CoreTestSuite struct {
	suite.Suite
}

func TestCoreTestSuite(t *testing.T) {
	suite.Run(t, new(CoreTestSuite))
}

func (s *CoreTestSuite) TestReportDelayWorksCorrectly() {
	dbpoolM := getDBPoolMock()
	loggerM := getLoggerMock()
	querierM := getQuerierMock()
	tripsHandlerM := getTripsHandlerMock(true)
	dtEstimatorM := getDeliveryTimeEstimatorMock()
	dqM := getDelayCheckQueueMock()
	nowFuncM := getNowFuncMock()

	c := core.New(
		dbpoolM,
		loggerM,
		querierM,
		tripsHandlerM,
		dtEstimatorM,
		dqM,
		nowFuncM,
	)

	req := &orderpb.AssignDelayReportToAgentRequest{
		ReportId: 1,
		AgentId:  1,
	}

	_, err := c.AssignDelayReportToAgent(context.Background(), req)
	s.NoError(err)
}

func getDBPoolMock() core.PgxPool {
	return nil
}

func getLoggerMock() core.Logger {
	return logrus.New()
}

func getQuerierMock() querier.Querier {
	m := &querierMocks.Querier{}
	m.On("CreateDelayCheck", mock.Anything, mock.Anything).Return(nil)
	// TODO: add more handlers

	return m
}

func getTripsHandlerMock(isOngoing bool) core.TripsHandler {
	m := &coreMocks.TripsHandler{}
	m.On("IsOngoing", mock.Anything).Return(isOngoing)

	return m
}

func getDeliveryTimeEstimatorMock() core.DeliveryTimeEstimator {
	m := &coreMocks.DeliveryTimeEstimator{}
	m.On("Estimate", mock.Anything, mock.Anything).Return(time.Hour)

	return m
}

func getDelayCheckQueueMock() core.DelayCheckQueue {
	m := &coreMocks.DelayCheckQueue{}
	m.On("Put", mock.Anything, mock.Anything).Return(nil)

	return m
}

func getNowFuncMock() func() time.Time {
	return time.Now
}
