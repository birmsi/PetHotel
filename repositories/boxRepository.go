package repositories

import (
	"PetHotel/models"
	"database/sql"
	"log"
	"log/slog"
	"time"
)

type BoxRepository struct {
	DB      *sql.DB
	slogger *slog.Logger
}

func NewBoxRepository(db *sql.DB, slogger *slog.Logger) BoxRepository {
	return BoxRepository{DB: db, slogger: slogger}
}

func (bx BoxRepository) CreateBox(box models.Box) (*int, error) {
	var lastInsertID int
	err := bx.DB.QueryRow("INSERT INTO boxes (number, size) VALUES ($1, $2) RETURNING id", box.Number, box.Size).Scan(&lastInsertID)

	if err != nil {
		return nil, err
	}

	return &lastInsertID, nil
}

func (bx BoxRepository) GetBox(id int) (*models.Box, error) {

	row := bx.DB.QueryRow("SELECT id, number, size FROM boxes WHERE id=$1", id)

	var box models.Box
	if err := row.Scan(&box.ID, &box.Number, &box.Size); err != nil {
		return nil, err
	}

	bx.slogger.Info("Get -> ", "id", box.ID, "number", box.Number, "size", box.Size)

	return &box, nil
}

func (bx BoxRepository) GetFutureAvailabilities(boxID int) ([]*models.Availability, error) {
	stmt := "SELECT start_time, end_time, price FROM availabilities WHERE box_id = $1 AND (start_time >= NOW() OR end_time >= NOW())"
	rows, err := bx.DB.Query(stmt, boxID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	availabilities := make([]*models.Availability, 0)
	var availability models.Availability

	for rows.Next() {
		if err = rows.Scan(&availability.StartTime, &availability.EndTime, &availability.Price); err != nil {
			return nil, err
		}
		availabilities = append(availabilities, &availability)
	}

	return availabilities, err
}

func (bx BoxRepository) AddAvailabilities(availabilities []models.Availability) error {

	transaction, err := bx.DB.Begin()
	if err != nil {
		return err
	}
	defer transaction.Rollback()

	stmt, err := transaction.Prepare(`
        INSERT INTO availabilities (box_id, start_time, end_time, price)
        VALUES ($1, $2, $3, $4)
    `)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, a := range availabilities {
		_, err := stmt.Exec(a.BoxID, a.StartTime, a.EndTime, a.Price)
		if err != nil {
			return err
		}
	}

	if err := transaction.Commit(); err != nil {
		return err
	}

	return nil
}

func (bx BoxRepository) GetAvailabilities(boxID int, start time.Time, end time.Time) ([]*models.Availability, error) {
	query := `
        SELECT id, box_id, start_time, end_time, price
        FROM availabilities
        WHERE box_id = $1
          AND (
            (start_time <= $2 AND end_time >= $3)
            OR (start_time >= $4 AND start_time <= $5)
          )
    `

	rows, err := bx.DB.Query(query, boxID, end, start, start, end)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	availabilities := make([]*models.Availability, 0)
	for rows.Next() {
		var availability models.Availability

		err := rows.Scan(&availability.ID, &availability.BoxID, &availability.StartTime, &availability.EndTime, &availability.Price)
		if err != nil {
			bx.slogger.Error("Failed to scan rows", err)
		}
		availabilities = append(availabilities, &availability)
	}

	return availabilities, nil
}
