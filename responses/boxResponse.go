package responses

type BoxResponse struct {
	ID            int
	Number        int
	Size          string
	Availabilites []*AvailabilityResponse
	Bookings      []*BookingResponse
}
