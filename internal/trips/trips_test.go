package trips_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/HomayoonAlimohammadi/order-delay/internal/database/querier"
	"github.com/HomayoonAlimohammadi/order-delay/internal/trips"
)

type TripsTestSuite struct {
	suite.Suite
}

func TestTripsTestSuite(t *testing.T) {
	suite.Run(t, new(TripsTestSuite))
}

func (s *TripsTestSuite) TestOngoingTrips() {
	h := trips.NewHandler()
	s.True(h.IsOngoing(querier.Trip{Status: querier.TripStatusAssigned}))
	s.True(h.IsOngoing(querier.Trip{Status: querier.TripStatusAtVendor}))
	s.True(h.IsOngoing(querier.Trip{Status: querier.TripStatusPicked}))
	s.False(h.IsOngoing(querier.Trip{Status: querier.TripStatusDelivered}))
}
