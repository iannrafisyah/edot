package controller

import (
	"context"
	"errors"
	"net/http"

	"Edot/models"
	"Edot/modules/stock/dto"
	"Edot/modules/stock/repository"
	"Edot/packages/logger"
	"Edot/utilities"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

type (
	IStockController interface {
		Create(ctx context.Context, reqData *dto.AddCartRequest, tx *gorm.DB) error
	}

	StockController struct {
		fx.In
		Logger          *logger.Logger
		StockRepository repository.IStockInterface
	}
)

// NewController :
func NewController(stockController StockController) IStockController {
	return &stockController
}

// Create :
func (r *StockController) Create(ctx context.Context, reqData *dto.AddCartRequest, tx *gorm.DB) error {
	// Validate request data
	messages, err := reqData.Validate()
	if err != nil {
		r.Logger.Error(err)
		return utilities.ErrorRequest(err, http.StatusBadRequest, messages)
	}

	fetchLatestStock, err := r.StockRepository.FindLatestStock(ctx, &models.Stock{
		WarehouseID: reqData.WarehouseID,
		ProductID:   reqData.ProductID,
	}, tx)
	if err != nil && err != gorm.ErrRecordNotFound {
		r.Logger.Error(err)
		return utilities.ErrorRequest(
			errors.New(utilities.InternalServiceError),
			http.StatusInternalServerError,
		)
	}

	newStock := models.Stock{
		Operator:      reqData.Operator,
		WarehouseID:   reqData.WarehouseID,
		ProductID:     reqData.ProductID,
		Qty:           reqData.Qty,
		Latest:        reqData.Qty,
		TransactionID: &reqData.TransactionID,
	}

	if fetchLatestStock != nil && err != gorm.ErrRecordNotFound {
		newStock.Previous = fetchLatestStock.Latest

		if reqData.Operator == models.StockOperatorDecrement && fetchLatestStock.Latest >= 1 {
			newStock.Latest = fetchLatestStock.Latest - reqData.Qty
		} else if reqData.Operator == models.StockOperatorIncrement {
			newStock.Latest = fetchLatestStock.Latest + reqData.Qty
		}
	}

	if _, err := r.StockRepository.Create(ctx, &newStock, tx); err != nil {
		r.Logger.Error(err)
		return utilities.ErrorRequest(
			errors.New(utilities.InternalServiceError),
			http.StatusInternalServerError,
		)
	}

	return nil
}
