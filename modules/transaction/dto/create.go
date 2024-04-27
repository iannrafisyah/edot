package dto

import (
	"errors"

	"Edot/models"
	"Edot/packages/validation"
	"Edot/utilities"
)

type (
	CreateRequest struct {
		ContextUserID int

		TransactionType models.TransactionType

		Validation validation.Validation
	}
)

func (r CreateRequest) Validate() (map[string]interface{}, error) {
	r.Validation.IsIntegerMin(r.ContextUserID, 1, "ContextUserID")

	if err := r.TransactionType.IsValid(); err != nil {
		r.Validation.NewError("TransactionType", err)
	}

	if len(r.Validation.Messages) > 0 {
		return r.Validation.Messages, errors.New(utilities.BadRequest)
	}

	return nil, nil
}
