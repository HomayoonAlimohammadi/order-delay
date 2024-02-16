package clients

import (
	"context"
	"fmt"
	"time"

	"github.com/HomayoonAlimohammadi/order-delay/internal/database/querier"
	orderpb "github.com/HomayoonAlimohammadi/order-delay/proto"
)

type delayCheckQueue struct {
	timeout time.Duration
}

func NewDelayCheckQueue(timeout time.Duration) *delayCheckQueue {
	return &delayCheckQueue{timeout: timeout}
}

func (e *delayCheckQueue) Put(ctx context.Context, ordr querier.Order) error {
	// TODO: implement some kind of queue, e.g. kafka
	ordrProto := &orderpb.Order{
		Id:                  ordr.ID,
		VendorId:            ordr.VendorID,
		DeliveryTimeMinutes: int32(ordr.DeliveryTime.Microseconds / 60_000_000), // 1min = 60 * 1_000_000 microseconds
	}

	fmt.Printf("putting order in the queue...: %+v", ordrProto)

	return nil
}
