package dto

import (
	"Edot/packages/validation"
)

type (
	FindAllRequest struct {
		ContextUserID int

		Limit int
		Page  int

		Validation validation.Validation
	}

	FindAllResponse struct {
		ID    int     `json:"id"`
		Name  string  `json:"name"`
		Price float64 `json:"price"`
	}
)

func (r FindAllRequest) Validate() (map[string]interface{}, error) {
	//r.Validation.IsIntegerMin(r.ID, 1, "Role")

	//if len(r.Validation.Messages) > 0 {
	//	return r.Validation.Messages, errors.New(utilities.BadRequest)
	//}

	return nil, nil
}
