// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package querier

import (
	"context"
)

type Querier interface {
	// Create a new delay check
	CreateDelayCheck(ctx context.Context, arg CreateDelayCheckParams) error
	// Report a new delay for an order
	CreateDelayReport(ctx context.Context, orderID int64) error
	// Get an order by ID
	GetOrderByID(ctx context.Context, id int64) (Order, error)
	// Get a trip record for a specific order_id
	GetTripByOrderID(ctx context.Context, orderID int64) (Trip, error)
	// Get all delay reports for a given vendor ID
	ListDelayReportsForVendor(ctx context.Context, arg ListDelayReportsForVendorParams) ([]DelayReport, error)
	// Updates the delivery time for a specific order by its ID
	UpdateOrderDeliveryTime(ctx context.Context, arg UpdateOrderDeliveryTimeParams) error
}

var _ Querier = (*Queries)(nil)