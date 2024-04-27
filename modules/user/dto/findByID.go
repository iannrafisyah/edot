package dto

import (
	"time"

	"Edot/packages/validation"
)

type (
	FindByIDRequest struct {
		ContextUserID int

		Validation validation.Validation
	}

	FindByIDResponse struct {
		ID        int       `json:"id"`
		FullName  string    `json:"full_name"`
		Email     string    `json:"email"`
		CreatedAt time.Time `json:"created_at"`
	}
)
