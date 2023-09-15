package handlers

import (
	"PetHotel/services"
	"log/slog"
)

type BookingHandlers struct {
	BookingService services.BookingService
	Slogger        *slog.Logger
}

func NewBookingHandlers(bookingService services.BookingService, slogger *slog.Logger) BookingHandlers {
	return BookingHandlers{
		BookingService: bookingService,
		Slogger:        slogger}
}
