package responses

type BookingResponse struct {
	ID       int
	BoxID    int
	UserID   int
	CheckIn  string
	CheckOut string
	Status   string
	Price    float64
}
