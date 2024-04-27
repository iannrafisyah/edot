package repository

import (
	"context"

	"Edot/models"

	"gorm.io/gorm"
)

// CreateTransactionProduct :
func (r *TransactionRepository) CreateTransactionProduct(ctx context.Context, reqData *models.TransactionProduct, tx *gorm.DB) (*int, error) {
	if err := tx.WithContext(ctx).Create(&reqData).Error; err != nil {
		r.Logger.Error(err)
		return nil, err
	}

	return &reqData.ID, nil
}
