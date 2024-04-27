package dto

import (
	"errors"
	"time"

	"Edot/packages/validation"
	"Edot/utilities"
)

type (
	FindByEmailRequest struct {
		ContextUserID int

		Email string

		Validation validation.Validation
	}

	FindByEmailResponse struct {
		ID        int       `json:"id"`
		FullName  string    `json:"full_name"`
		Email     string    `json:"email"`
		Password  string    `json:"password"`
		CreatedAt time.Time `json:"created_at"`
	}
)

func (r FindByEmailRequest) Validate() (map[string]interface{}, error) {
	r.Validation.IsEmptyString(r.Email, "Email")
	r.Validation.IsEmailValid(r.Email, "Email")

	if len(r.Validation.Messages) > 0 {
		return r.Validation.Messages, errors.New(utilities.BadRequest)
	}

	return nil, nil
}
