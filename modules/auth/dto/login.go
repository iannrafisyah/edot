package dto

import (
	"errors"

	"Edot/packages/validation"
	"Edot/utilities"
)

type (
	LoginRequest struct {
		ContextUserID int

		Email    string
		Password string

		Validation validation.Validation
	}

	LoginResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
)

func (r LoginRequest) Validate() (map[string]interface{}, error) {
	r.Validation.IsEmptyString(r.Password, "Password")
	r.Validation.IsEmptyString(r.Email, "Email")
	r.Validation.IsEmailValid(r.Email, "Email")
	if len(r.Validation.Messages) > 0 {
		return r.Validation.Messages, errors.New(utilities.BadRequest)
	}

	return nil, nil
}
