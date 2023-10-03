package services

import (
	"PetHotel/models"
	"PetHotel/repositories"
	"log/slog"
	"time"
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
func (bx BoxService) GetBoxes() ([]*models.Box, error) {
	return bx.Repository.GetBoxes()
}

func (bx BoxService) GetFutureAvailabilities(boxID int) ([]*models.Availability, error) {
	return bx.Repository.GetFutureAvailabilities(boxID)
}

func (bx BoxService) AddAvailabilities(availabilities []models.Availability) error {
	return bx.Repository.AddAvailabilities(availabilities)
}

func (bx BoxService) GetAvailabilities(boxID int, start time.Time, end time.Time) ([]*models.Availability, error) {
	return bx.Repository.GetAvailabilities(boxID, start, end)
}
