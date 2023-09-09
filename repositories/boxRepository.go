package repositories

import (
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
