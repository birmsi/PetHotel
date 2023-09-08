package repositories

import "database/sql"

type BoxRepository struct {
	DB *sql.DB
}

func NewBoxRepository(db *sql.DB) BoxRepository {
	return BoxRepository{DB: db}
}
