package dto

import (
	"errors"

	"Edot/packages/validation"
	"Edot/utilities"
)

type (
	PaymentRequest struct {
		ContextUserID int

		ID int

		Validation validation.Validation
	}
)

func (r PaymentRequest) Validate() (map[string]interface{}, error) {
	r.Validation.IsIntegerMin(r.ID, 1, "ID")

	if len(r.Validation.Messages) > 0 {
		return r.Validation.Messages, errors.New(utilities.BadRequest)
	}

	return nil, nil
}
