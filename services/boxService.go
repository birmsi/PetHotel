package services

import (
	"PetHotel/models"
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

func (bx BoxService) CreateBox(box models.Box) (*int, error) {
	return bx.Repository.CreateBox(box)
}

func (bx BoxService) GetBox(id int) (*models.Box, error) {
	return bx.Repository.GetBox(id)
}
