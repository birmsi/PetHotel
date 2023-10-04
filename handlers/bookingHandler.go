package handlers

import (
	"PetHotel/services"
	"log/slog"
	"net/http"
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

func (bh BookingHandlers) GetAllBookings(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)

}
func (bh BookingHandlers) GetBooking(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (bh BookingHandlers) CreateBookingView(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}
func (bh BookingHandlers) CreateBookingPost(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}
func (bh BookingHandlers) UpdateBookingView(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}
func (bh BookingHandlers) UpdateBookingPut(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}
func (bh BookingHandlers) DeleteBookingPost(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}
