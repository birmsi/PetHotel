package models

import "time"

type OperationControl struct {
	CreatedAT     time.Time
	CreatedBy     int
	LastUpdatedAT time.Time
	LastUpdatedBy int
	DeletedAT     time.Time
	DeletedBy     int
}
