package repositories

import (
	"PetHotel/models"
	"database/sql"
	"log/slog"
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
