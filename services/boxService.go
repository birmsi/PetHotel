package services

import (
	"PetHotel/repositories"
	"log/slog"
)

type BoxService struct {
	Repository repositories.BoxRepository
	slogger    *slog.Logger
}

func NewService(repository repositories.BoxRepository, slogger *slog.Logger) BoxService {
	return BoxService{Repository: repository, slogger: slogger}
}
