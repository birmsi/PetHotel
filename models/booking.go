package models

import "time"

type Booking struct {
	Box      Box
	CheckIn  time.Time
	CheckOut time.Time
}
