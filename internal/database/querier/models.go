// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package querier

import (
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type TripStatus string

const (
	TripStatusAssigned  TripStatus = "assigned"
	TripStatusAtVendor  TripStatus = "at_vendor"
	TripStatusPicked    TripStatus = "picked"
	TripStatusDelivered TripStatus = "delivered"
)

func (e *TripStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = TripStatus(s)
	case string:
		*e = TripStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for TripStatus: %T", src)
	}
	return nil
}

type NullTripStatus struct {
	TripStatus TripStatus
	Valid      bool // Valid is true if TripStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullTripStatus) Scan(value interface{}) error {
	if value == nil {
		ns.TripStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.TripStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullTripStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.TripStatus), nil
}

type Agent struct {
	ID   int32
	Name string
}

type DelayCheck struct {
	ID       int64
	AgentID  int32
	ReportID int64
}

type DelayReport struct {
	ID        int64
	CreatedAt pgtype.Timestamptz
	OrderID   int64
}

type Order struct {
	ID           int64
	VendorID     int32
	DeliveryTime pgtype.Interval
	CreatedAt    pgtype.Timestamptz
}

type Trip struct {
	ID        int64
	Status    TripStatus
	OrderID   int64
	CreatedAt pgtype.Timestamptz
}

type Vendor struct {
	ID   int32
	Name string
}