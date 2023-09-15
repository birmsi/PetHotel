package models

import "time"

type Booking struct {
	OperationControl
	ID       int
	BoxID    int
	UserID   int
	CheckIn  time.Time
	CheckOut time.Time
	Status   string
	Price    float64
}
