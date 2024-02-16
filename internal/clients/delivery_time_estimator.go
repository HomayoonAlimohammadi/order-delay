package clients

import (
	"context"
	"time"

	"github.com/HomayoonAlimohammadi/order-delay/internal/database/querier"
)

type deliveryTimeEstimator struct {
	timeout time.Duration
}

func NewDeliveryTimeEstimator(timeout time.Duration) *deliveryTimeEstimator {
	return &deliveryTimeEstimator{timeout: timeout}
}

func (e *deliveryTimeEstimator) Estimate(ctx context.Context, order querier.Order) (time.Duration, error) {
	// TODO: implement the logic to call another service
	return time.Hour, nil
}
