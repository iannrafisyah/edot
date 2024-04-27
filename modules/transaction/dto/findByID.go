package dto

import (
	"errors"
	"time"

	"Edot/packages/validation"
	"Edot/utilities"
)

type (
	FindByIDRequest struct {
		ContextUserID int

		ID int

		Validation validation.Validation
	}

	FindByIDResponse struct {
		ID         int                `json:"id"`
		Invoice    string             `json:"invoice"`
		Total      float64            `json:"total"`
		Tax        float64            `json:"tax"`
		GrandTotal float64            `json:"grand_total"`
		Status     string             `json:"status"`
		Type       string             `json:"type"`
		Products   []*FindByIDProduct `json:"products"`
		CreatedAt  time.Time          `json:"created_at"`
	}

	FindByIDProduct struct {
		Name  string  `json:"name"`
		Price float64 `json:"price"`
		Qty   int     `json:"qty"`
		Total float64 `json:"total"`
	}
)

func (r FindByIDRequest) Validate() (map[string]interface{}, error) {
	r.Validation.IsIntegerMin(r.ID, 1, "TransactionID")

	if len(r.Validation.Messages) > 0 {
		return r.Validation.Messages, errors.New(utilities.BadRequest)
	}

	return nil, nil
}
