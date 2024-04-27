package repository

import (
	"context"

	"Edot/models"
	"Edot/packages/logger"
	"Edot/packages/postgres"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

type (
	IStockInterface interface {
		Create(ctx context.Context, reqData *models.Stock, tx *gorm.DB) (*int, error)
		FindLatestStock(ctx context.Context, reqData *models.Stock, tx *gorm.DB) (*models.Stock, error)
	}

	StockRepository struct {
		fx.In
		DB     *postgres.DB
		Logger *logger.Logger
	}
)

// NewRepository :
func NewRepository(stockRepository StockRepository) IStockInterface {
	return &stockRepository
}

// Create :
func (r *StockRepository) Create(ctx context.Context, reqData *models.Stock, tx *gorm.DB) (*int, error) {
	if err := tx.WithContext(ctx).Create(&reqData).Error; err != nil {
		r.Logger.Error(err)
		return nil, err
	}

	return &reqData.ID, nil
}

// FindLatestStock :
func (r *StockRepository) FindLatestStock(ctx context.Context, reqData *models.Stock, tx *gorm.DB) (*models.Stock, error) {
	stock := new(models.Stock)

	if err := tx.WithContext(ctx).
		Where("warehouse_id", reqData.WarehouseID).
		Where("product_id", reqData.ProductID).
		Last(&stock).Error; err != nil {
		r.Logger.Error(err)
		return nil, err
	}

	return stock, nil
}
