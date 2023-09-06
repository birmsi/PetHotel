package models

import "time"

type Service interface {
	Book()
	IsAvailable(start, end time.Time) bool
	GetPrice()
}
