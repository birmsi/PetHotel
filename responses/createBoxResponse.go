package responses

import "PetHotel/models"

type CreateBoxViewResponse struct {
	ErrorMessage string
	BoxSizes     []string
	Box          models.Box
}
