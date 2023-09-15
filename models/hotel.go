package models

import (
	"PetHotel/helpers"
	"fmt"
	"time"
)

type HotelService struct {
	Boxes    []Box
	Bookings []Booking
}

type Box struct {
	OperationControl
	ID              int
	Number          int
	Size            Size
	AvailabilityIDs []int
}

type Size string

const SmallSize Size = "small"
const MediumSize Size = "medium"
const LargeSize Size = "large"

type Availability struct {
	ID        int
	BoxID     int
	StartTime time.Time
	EndTime   time.Time
	Price     float64
}

func (h *HotelService) Book(box Box, checkIn, checkOut time.Time, availabilities []Availability) string {

	if !h.IsAvailable(box, checkIn, checkOut, availabilities) {
		return fmt.Sprintf("Box %d (%s) is not available for the selected time period.\n", box.Number, box.Size)
	}

	totalPrice := 0.0
	for _, dateRange := range availabilities {
		if dateRange.StartTime.Before(checkOut) && dateRange.EndTime.After(checkIn) {

			start := helpers.Max(dateRange.StartTime, checkIn)
			end := helpers.Min(dateRange.EndTime, checkOut)
			totalPrice += dateRange.Price * end.Sub(start).Hours() / 24.0
		}
	}

	booking := Booking{
		BoxID:    box.ID,
		CheckIn:  checkIn,
		CheckOut: checkOut,
	}

	h.Bookings = append(h.Bookings, booking)

	return fmt.Sprintf("Booked from %s to %s in Box %d (%s) for $%.2f",
		checkIn.Format("2006-01-02"), checkOut.Format("2006-01-02"), box.Number, box.Size, totalPrice)
}

func (h *HotelService) IsAvailable(box Box, checkIn, checkOut time.Time, availabilities []Availability) bool {
	if len(availabilities) == 0 {
		//maybe log this as it should have some availability
		return false
	}
	isAvailable := false
	for _, dateRange := range availabilities {
		if dateRange.StartTime.Before(checkIn) && dateRange.EndTime.After(checkOut) {
			isAvailable = true
		}
		if dateRange.StartTime.Before(checkIn) && dateRange.EndTime.Equal(checkOut) {
			isAvailable = true
		}
		if dateRange.StartTime.Equal(checkIn) && dateRange.EndTime.After(checkOut) {
			isAvailable = true
		}
		if dateRange.StartTime.Equal(checkIn) && dateRange.EndTime.Equal(checkOut) {
			isAvailable = true
		}
	}

	if !isAvailable {
		return false
	}

	canBook := true
	//Add tests to this logic :p
	for _, booking := range h.Bookings {

		if booking.BoxID == box.ID &&
			(booking.CheckIn.Before(checkIn) && booking.CheckOut.After(checkIn) ||
				booking.CheckIn.Equal(checkIn) ||
				booking.CheckIn.After(checkIn) && booking.CheckOut.Before(checkOut) ||
				booking.CheckIn.After(checkIn) && booking.CheckOut.Equal(checkOut) ||
				booking.CheckIn.After(checkIn) && booking.CheckIn.Before(checkOut) && booking.CheckOut.After(checkOut)) {

			canBook = false
		}
	}

	return canBook
}
