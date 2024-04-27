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
	ICartInterface interface {
		Create(ctx context.Context, reqData *models.Cart, tx *gorm.DB) (*int, error)
		FindByUserID(ctx context.Context, reqData *models.Cart, tx *gorm.DB) ([]*models.Cart, error)
		Delete(ctx context.Context, reqData *models.Cart, tx *gorm.DB) error
		TotalQtyWithTrxPending(ctx context.Context, reqData *models.Cart, tx *gorm.DB) (*int64, error)
	}

	CartRepository struct {
		fx.In
		DB     *postgres.DB
		Logger *logger.Logger
	}
)

// NewRepository :
func NewRepository(cartRepository CartRepository) ICartInterface {
	return &cartRepository
}

// Create :
func (r *CartRepository) Create(ctx context.Context, reqData *models.Cart, tx *gorm.DB) (*int, error) {
	if err := tx.WithContext(ctx).Create(&reqData).Error; err != nil {
		r.Logger.Error(err)
		return nil, err
	}

	return &reqData.ID, nil
}

// Delete :
func (r *CartRepository) Delete(ctx context.Context, reqData *models.Cart, tx *gorm.DB) error {
	if err := tx.WithContext(ctx).Where(
		"user_id", reqData.UserID,
	).Delete(&models.Cart{}).Error; err != nil {
		r.Logger.Error(err)
		return err
	}

	return nil
}

// FindByUserID :
func (r *CartRepository) FindByUserID(ctx context.Context, reqData *models.Cart, tx *gorm.DB) ([]*models.Cart, error) {
	carts := make([]*models.Cart, 0)

	if err := tx.WithContext(ctx).
		Where(&models.Cart{
			UserID: reqData.UserID,
		}).
		Find(&carts).Error; err != nil {
		r.Logger.Error(err)
		return nil, err
	}

	return carts, nil
}

// TotalQtyWithTrxPending :
func (r *CartRepository) TotalQtyWithTrxPending(ctx context.Context, reqData *models.Cart, tx *gorm.DB) (*int64, error) {
	totalQty := int64(0)

	if err := tx.WithContext(ctx).Select("coalesce(sum(transaction_products.qty),0) as total_qty").Model(&models.TransactionProduct{}).
		Joins("inner join transactions trx on trx.id = transaction_products.transaction_id").
		Where("product_id", reqData.ProductID).
		Where("trx.status", models.TransactionStatusUnpaid).
		Where("warehouse_id", reqData.WarehouseID).
		Take(&totalQty).Error; err != nil {
		r.Logger.Error(err)
		return nil, err
	}

	return &totalQty, nil
}
