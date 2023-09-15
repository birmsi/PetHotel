package services

import (
	"PetHotel/repositories"
	"log/slog"
	"time"

	"PetHotel/models"
)

type BookingService struct {
	BookingRepository repositories.BookingRepository
	Slogger           *slog.Logger
}

func NewBookingService(bookingRepository repositories.BookingRepository, slogger *slog.Logger) BookingService {
	return BookingService{
		BookingRepository: bookingRepository,
		Slogger:           slogger,
	}
}

func (bs BookingService) GetFutureBookings(boxID int) ([]*models.Booking, error) {
	return bs.BookingRepository.GetFutureBookings(boxID)
}

func (bs BookingService) GetBookings(boxID int, start time.Time, end time.Time) ([]*models.Booking, error) {
	return bs.BookingRepository.GetBookings(boxID, start, end)
}
