package services

import "PetHotel/repositories"

type BoxService struct {
	Repository repositories.BoxRepository
}

func NewService(repository repositories.BoxRepository) BoxService {
	return BoxService{Repository: repository}
}
