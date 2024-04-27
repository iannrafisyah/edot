package dto

import (
	"Edot/packages/validation"
	"Edot/utilities"
	"errors"
)

type (
	DeleteRequest struct {
		ContextUserID int

		Validation validation.Validation
	}
)

func (r DeleteRequest) Validate() (map[string]interface{}, error) {
	r.Validation.IsIntegerMin(r.ContextUserID, 1, "ContextUserID")

	if len(r.Validation.Messages) > 0 {
		return r.Validation.Messages, errors.New(utilities.BadRequest)
	}

	return nil, nil
}
