package dto

import (
	"Edot/packages/validation"
	"Edot/utilities"
	"errors"
)

type (
	FindAllRequest struct {
		ContextUserID int
		Validation    validation.Validation
	}

	FindAllResponse struct {
		ProductID     int     `json:"product_id"`
		Name          string  `json:"name"`
		Price         float64 `json:"price"`
		Qty           int     `json:"qty"`
		WarehouseID   int     `json:"warehouse_id"`
		ToWarehouseID *int    `json:"to_warehouse_id"`
	}
)

func (r FindAllRequest) Validate() (map[string]interface{}, error) {
	r.Validation.IsIntegerMin(r.ContextUserID, 1, "ContextUserID")

	if len(r.Validation.Messages) > 0 {
		return r.Validation.Messages, errors.New(utilities.BadRequest)
	}

	return nil, nil
}
