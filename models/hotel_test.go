package models

import (
	"testing"
	"time"
)

func TestNoAvailableTimes(t *testing.T) {
	now := time.Now()
	hotel := HotelService{}
	testBox := Box{Number: 123, Size: LargeSize, Availability: nil}

	result := hotel.IsAvailable(
		testBox,
		time.Now().Add(24*time.Hour),
		now.Add(72*time.Hour))

	if result == true {
		t.Errorf("Test no available dates empty - failed. Expected false, got %v", result)
	}
}

func TestIsAvailableNoAvailability(t *testing.T) {
	now := time.Now()
	hotel := HotelService{}
	testBox := Box{
		Number: 123,
		Size:   LargeSize,
		Availability: []DateRangeWithPrice{
			{
				Start: now,
				End:   now.Add(1 * time.Hour),
			}},
	}

	result := hotel.IsAvailable(
		testBox,
		now.Add(24*time.Hour),
		now.Add(48*time.Hour))

	if result == true {
		t.Errorf("Test available dates before requested - failed. Expected false, got %v", result)
	}

	testBox.Availability = append(testBox.Availability, DateRangeWithPrice{
		Start: now.Add(72 * time.Hour),
		End:   now.Add(96 * time.Hour),
	})

	result = hotel.IsAvailable(
		testBox,
		now.Add(24*time.Hour),
		now.Add(48*time.Hour))

	if result == true {
		t.Errorf("Test available dates before requested - failed. Expected false, got %v", result)
	}
}

func TestIsAvailable(t *testing.T) {
	now := time.Now()
	hotel := HotelService{}
	testBox := Box{
		Number: 123,
		Size:   LargeSize,
		Availability: []DateRangeWithPrice{
			{
				Start: now,
				End:   now.Add(100 * time.Hour),
			}},
	}

	result := hotel.IsAvailable(
		testBox,
		now.Add(24*time.Hour),
		now.Add(48*time.Hour))

	if result == false {
		t.Errorf("Test available dates - failed. Expected true, got %v", result)
	}
}

func TestIsAvailableWithBookingFail(t *testing.T) {
	now := time.Now()
	box := Box{
		Number: 123,
		Size:   LargeSize,
		Availability: []DateRangeWithPrice{
			{
				Start: now,
				End:   now.Add(100 * time.Hour),
			}},
	}
	hotel := HotelService{
		Boxes: []Box{box},
		Bookings: []Booking{{
			Box:      Box{Number: 123},
			CheckIn:  now.Add(24 * time.Hour),
			CheckOut: now.Add(72 * time.Hour),
		}},
	}
	result := hotel.IsAvailable(
		box,
		now,
		now.Add(30*time.Hour))

	if result == true {
		t.Errorf("test failed - should have no availability due to booking")
	}

	result = hotel.IsAvailable(
		box,
		now.Add(48*time.Hour),
		now.Add(96*time.Hour))

	if result == true {
		t.Errorf("test failed - should have no availability due to booking")
	}

	result = hotel.IsAvailable(
		box,
		now.Add(48*time.Hour),
		now.Add(80*time.Hour))

	if result == true {
		t.Errorf("test failed - beginning at the mid of the booking and ending after - should have no availability due to booking")
	}

	result = hotel.IsAvailable(
		box,
		now.Add(24*time.Hour),
		now.Add(72*time.Hour))

	if result == true {
		t.Errorf("test failed - exactly same time as booking - should have no availability due to booking")
	}

}

func TestIsAvailableWithBookingPass(t *testing.T) {
	now := time.Now()
	box := Box{
		Number: 123,
		Size:   LargeSize,
		Availability: []DateRangeWithPrice{
			{
				Start: now,
				End:   now.Add(100 * time.Hour),
			}},
	}
	hotel := HotelService{
		Boxes: []Box{box},
		Bookings: []Booking{{
			Box:      Box{Number: 123},
			CheckIn:  now.Add(24 * time.Hour),
			CheckOut: now.Add(72 * time.Hour),
		}},
	}

	result := hotel.IsAvailable(
		box,
		now.Add(72*time.Hour),
		now.Add(96*time.Hour))

	if result == false {
		t.Errorf("test failed - should have no availability due to booking")
	}
}

func TestBook(t *testing.T) {

}
