package dto

import (
	"Edot/models"
	"Edot/packages/validation"
	"Edot/utilities"
	"errors"
)

type (
	AddCartRequest struct {
		ContextUserID int

		Qty           int
		Operator      models.StockOperator
		WarehouseID   int
		ProductID     int
		TransactionID int

		Validation validation.Validation
	}
)

func (r AddCartRequest) Validate() (map[string]interface{}, error) {
	r.Validation.IsIntegerMin(r.ContextUserID, 1, "ContextUserID")
	r.Validation.IsIntegerMin(r.WarehouseID, 1, "WarehouseID")
	r.Validation.IsIntegerMin(r.TransactionID, 1, "TransactionID")
	r.Validation.IsIntegerMin(r.ProductID, 1, "ProductID")
	r.Validation.IsIntegerMin(r.Qty, 1, "Qty")

	if len(r.Validation.Messages) > 0 {
		return r.Validation.Messages, errors.New(utilities.BadRequest)
	}

	return nil, nil
}
