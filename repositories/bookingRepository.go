package repositories

import (
	"PetHotel/models"
	"database/sql"
	"log"
	"log/slog"
	"time"
)

type BookingRepository struct {
	DB      *sql.DB
	Slogger *slog.Logger
}

func NewBookingRepository(db *sql.DB, slogger *slog.Logger) BookingRepository {
	return BookingRepository{
		DB:      db,
		Slogger: slogger,
	}
}

func (br BookingRepository) GetFutureBookings(boxID int) ([]*models.Booking, error) {
	stmt := "SELECT checkin, checkout FROM bookings WHERE box_id = $1 AND checkin >= NOW() OR checkin < NOW() AND checkout >= NOW()"
	rows, err := br.DB.Query(stmt, boxID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	futureBookings := make([]*models.Booking, 0)

	for rows.Next() {
		var futureBooking models.Booking
		rows.Scan(&futureBooking.CheckIn, futureBooking.CheckOut)

		futureBookings = append(futureBookings, &futureBooking)
	}

	return futureBookings, nil
}

func (br BookingRepository) GetBookings(boxID int, start time.Time, end time.Time) ([]*models.Booking, error) {
	query := `
	SELECT id, box_id, checkin, checkout
	FROM bookings
	WHERE box_id = $1
	  AND (
		(checkin <= $2 AND checkout >= $3)
		OR (checkin >= $4 AND checkin <= $5)
	  )
`

	rows, err := br.DB.Query(query, boxID, end, start, start, end)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	bookings := make([]*models.Booking, 0)
	for rows.Next() {
		var booking models.Booking

		err := rows.Scan(&booking.ID, &booking.BoxID, &booking.CheckIn, &booking.CheckOut)
		if err != nil {
			br.Slogger.Error("Failed to scan rows", err)
		}
		bookings = append(bookings, &booking)
	}

	return bookings, nil
}
