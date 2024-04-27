package dto

import (
	"Edot/packages/validation"
	"Edot/utilities"
	"errors"
)

type (
	AddCartRequest struct {
		ContextUserID int

		Items []AddCartItem

		Validation validation.Validation
	}

	AddCartItem struct {
		ProductID     int
		WarehouseID   int
		ToWarehouseID *int
		Qty           int
	}
)

func (r AddCartRequest) Validate() (map[string]interface{}, error) {
	r.Validation.IsIntegerMin(r.ContextUserID, 1, "ContextUserID")
	r.Validation.IsIntegerMin(len(r.Items), 1, "Items")

	for _, v := range r.Items {
		r.Validation.IsIntegerMin(v.ProductID, 1, "ProductID")
		r.Validation.IsIntegerMin(v.Qty, 1, "Qty")
		r.Validation.IsIntegerMin(v.WarehouseID, 1, "WarehouseID")
	}

	if len(r.Validation.Messages) > 0 {
		return r.Validation.Messages, errors.New(utilities.BadRequest)
	}

	return nil, nil
}
