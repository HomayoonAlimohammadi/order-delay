package trips

import (
	"slices"

	"github.com/HomayoonAlimohammadi/order-delay/internal/database/querier"
)

type handler struct{}

func NewHandler() *handler {
	return &handler{}
}

func (h *handler) IsOngoing(trip querier.Trip) bool {
	return slices.Contains([]querier.TripStatus{
		querier.TripStatusAssigned,
		querier.TripStatusAtVendor,
		querier.TripStatusPicked,
	}, trip.Status)
}
