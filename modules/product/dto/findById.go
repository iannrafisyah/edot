package dto

import (
	"errors"

	"Edot/packages/validation"
	"Edot/utilities"
)

type (
	FindByIDRequest struct {
		ContextUserID int

		ID          int
		WarehouseID int

		Validation validation.Validation
	}

	FindByIDResponse struct {
		ID             int                       `json:"id"`
		Name           string                    `json:"name"`
		Price          float64                   `json:"price"`
		StockAvailable int                       `json:"stock_available"`
		Warehouses     []*FindByIDWarehouseStock `json:"warehouses"`
	}

	FindByIDWarehouseStock struct {
		WarehouseName string `json:"warehouse_name"`
		WarehouseID   int    `json:"warehouse_id"`
		Latest        int    `json:"latest"`
		Previous      int    `json:"previous"`
		Qty           int    `json:"qty"`
	}
)

func (r FindByIDRequest) Validate() (map[string]interface{}, error) {
	r.Validation.IsIntegerMin(r.ID, 1, "ProductID")

	if len(r.Validation.Messages) > 0 {
		return r.Validation.Messages, errors.New(utilities.BadRequest)
	}

	return nil, nil
}
