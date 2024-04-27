package dto

import (
	"errors"
	"time"

	"Edot/models"
	"Edot/packages/validation"
	"Edot/utilities"
)

type (
	FindAllRequest struct {
		ContextUserID int

		Limit int
		Page  int
		Type  models.TransactionType

		Validation validation.Validation
	}

	FindAllResponse struct {
		ID         int       `json:"id"`
		Invoice    string    `json:"invoice"`
		Total      float64   `json:"total"`
		Tax        float64   `json:"tax"`
		GrandTotal float64   `json:"grand_total"`
		Status     string    `json:"status"`
		Type       string    `json:"type"`
		CreatedAt  time.Time `json:"created_at"`
	}
)

func (r FindAllRequest) Validate() (map[string]interface{}, error) {
	r.Validation.IsIntegerMin(r.ContextUserID, 1, "ContextUserID")

	if len(r.Validation.Messages) > 0 {
		return r.Validation.Messages, errors.New(utilities.BadRequest)
	}

	return nil, nil
}
