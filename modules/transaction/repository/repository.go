package repository

import (
	"context"
	"time"

	"Edot/models"
	"Edot/packages/logger"
	"Edot/packages/paginate"
	"Edot/packages/postgres"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

type (
	ITransactionInterface interface {
		FindByID(ctx context.Context, reqData *models.Transaction, tx *gorm.DB) (*models.Transaction, error)
		FindAll(ctx context.Context, reqData *models.Transaction, pagination paginate.Pagination, tx *gorm.DB) ([]*models.Transaction, *paginate.Pagination, error)
		Create(ctx context.Context, reqData *models.Transaction, tx *gorm.DB) (*int, error)
		UpdateStatus(ctx context.Context, reqData *models.Transaction, tx *gorm.DB) error
		FindAllOrderExpired(ctx context.Context, reqData *models.Transaction, tx *gorm.DB) ([]*models.Transaction, error)
		CreateTransactionProduct(ctx context.Context, reqData *models.TransactionProduct, tx *gorm.DB) (*int, error)
	}

	TransactionRepository struct {
		fx.In
		DB     *postgres.DB
		Logger *logger.Logger
	}
)

// NewRepository :
func NewRepository(transactionRepository TransactionRepository) ITransactionInterface {
	return &transactionRepository
}

// Create :
func (r *TransactionRepository) Create(ctx context.Context, reqData *models.Transaction, tx *gorm.DB) (*int, error) {
	if err := tx.WithContext(ctx).Create(&reqData).Error; err != nil {
		r.Logger.Error(err)
		return nil, err
	}

	return &reqData.ID, nil
}

// FindByID :
func (r *TransactionRepository) FindByID(ctx context.Context, reqData *models.Transaction, tx *gorm.DB) (*models.Transaction, error) {
	transaction := new(models.Transaction)

	if err := tx.WithContext(ctx).
		Preload("TransactionProducts").
		Where(&models.Transaction{
			ID:     reqData.ID,
			UserID: reqData.UserID,
		}).
		Take(&transaction).Error; err != nil {
		r.Logger.Error(err)
		return nil, err
	}

	return transaction, nil
}

// FindAll :
func (r *TransactionRepository) FindAll(ctx context.Context, reqData *models.Transaction, pagination paginate.Pagination, tx *gorm.DB) ([]*models.Transaction, *paginate.Pagination, error) {
	transactions := make([]*models.Transaction, 0)

	queryBuilder := tx.WithContext(ctx)

	if err := queryBuilder.Where("user_id", reqData.UserID).
		Where("type", reqData.Type).
		Order("id desc").Scopes(paginate.Paginate(reqData, &pagination, queryBuilder)).
		Find(&transactions).Error; err != nil {
		r.Logger.Error(err)
		return nil, nil, err
	}

	return transactions, &pagination, nil
}

// UpdateStatus :
func (r *TransactionRepository) UpdateStatus(ctx context.Context, reqData *models.Transaction, tx *gorm.DB) error {
	if err := tx.WithContext(ctx).
		Model(&models.Transaction{}).
		Where("id", reqData.ID).
		Where("user_id", reqData.UserID).
		Updates(map[string]interface{}{
			"status":     reqData.Status,
			"updated_at": time.Now(),
		}).
		Error; err != nil {
		r.Logger.Error(err)
		return err
	}
	return nil
}

// FindAllOrderExpired :
func (r *TransactionRepository) FindAllOrderExpired(ctx context.Context, reqData *models.Transaction, tx *gorm.DB) ([]*models.Transaction, error) {
	transactions := make([]*models.Transaction, 0)

	queryBuilder := tx.WithContext(ctx)

	if err := queryBuilder.
		Where("status", models.TransactionStatusUnpaid).
		Where("created_at <= now() - interval '1 minutes'").
		Order("id desc").
		Find(&transactions).Error; err != nil {
		r.Logger.Error(err)
		return nil, err
	}

	return transactions, nil
}
