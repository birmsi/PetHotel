package responses

type AvailabilityResponse struct {
	ID        int
	BoxID     int
	StartTime string
	EndTime   string
	Price     float64
}
