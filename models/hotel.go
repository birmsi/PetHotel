package models

import (
	"PetHotel/helpers"
	"fmt"
	"time"
)

type HotelService struct {
	Name     string
	Boxes    []Box
	Bookings []Booking
}

type Box struct {
	Number       int
	Size         Size
	Occupancy    int
	Availability []DateRangeWithPrice
}

type Size string

const SmallSize Size = "small"
const MediumSize Size = "medium"
const LargeSize Size = "large"

type DateRangeWithPrice struct {
	Start time.Time
	End   time.Time
	Price float64
}

func (h *HotelService) Book(box Box, checkIn, checkOut time.Time) string {

	if !h.IsAvailable(box, checkIn, checkOut) {
		return fmt.Sprintf("Box %d (%s) is not available for the selected time period.\n", box.Number, box.Size)
	}

	totalPrice := 0.0
	for _, dateRange := range box.Availability {
		if dateRange.Start.Before(checkOut) && dateRange.End.After(checkIn) {

			start := helpers.Max(dateRange.Start, checkIn)
			end := helpers.Min(dateRange.End, checkOut)
			totalPrice += dateRange.Price * end.Sub(start).Hours() / 24.0
		}
	}

	booking := Booking{
		Box:      box,
		CheckIn:  checkIn,
		CheckOut: checkOut,
	}

	h.Bookings = append(h.Bookings, booking)

	return fmt.Sprintf("Booked %s from %s to %s in Box %d (%s) for $%.2f",
		h.Name, checkIn.Format("2006-01-02"), checkOut.Format("2006-01-02"), box.Number, box.Size, totalPrice)
}

func (h *HotelService) IsAvailable(box Box, checkIn, checkOut time.Time) bool {
	for _, dateRange := range box.Availability {
		if dateRange.Start.Before(checkIn) && dateRange.End.After(checkOut) {
			return true
		}
		if dateRange.Start.Before(checkIn) && dateRange.End.Equal(checkOut) {
			return true
		}
		if dateRange.Start.Equal(checkIn) && dateRange.End.After(checkOut) {
			return true
		}
		if dateRange.Start.Equal(checkIn) && dateRange.End.Equal(checkOut) {
			return true
		}
	}

	canBook := true
	//Add tests to this logic :p
	for _, booking := range h.Bookings {
		if booking.Box.Number == box.Number && (booking.CheckIn.Before(checkIn) && booking.CheckOut.After(checkIn) ||
			booking.CheckIn.Equal(checkIn) ||
			booking.CheckIn.After(checkIn) && booking.CheckOut.Before(checkOut) ||
			booking.CheckIn.After(checkIn) && booking.CheckOut.Equal(checkOut)) {

			canBook = false
		}
	}

	return canBook
}

func (h *HotelService) AddBox(box Box) {
	h.Boxes = append(h.Boxes, box)
}
