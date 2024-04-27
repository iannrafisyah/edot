package repository

import (
	"context"

	"Edot/models"
	"Edot/packages/logger"
	"Edot/packages/paginate"
	"Edot/packages/postgres"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

type (
	IProductInterface interface {
		FindAll(ctx context.Context, reqData *models.Product, pagination paginate.Pagination, tx *gorm.DB) ([]*models.Product, *paginate.Pagination, error)
		Create(ctx context.Context, reqData *models.Product, tx *gorm.DB) (*int, error)
		FindByID(ctx context.Context, reqData *models.Product, tx *gorm.DB) (*models.Product, error)
	}

	ProductRepository struct {
		fx.In
		DB     *postgres.DB
		Logger *logger.Logger
	}
)

// NewRepository :
func NewRepository(productRepository ProductRepository) IProductInterface {
	return &productRepository
}

// Create :
func (r *ProductRepository) Create(ctx context.Context, reqData *models.Product, tx *gorm.DB) (*int, error) {
	if err := tx.WithContext(ctx).Create(&reqData).Error; err != nil {
		r.Logger.Error(err)
		return nil, err
	}

	return &reqData.ID, nil
}

// FindByID :
func (r *ProductRepository) FindByID(ctx context.Context, reqData *models.Product, tx *gorm.DB) (*models.Product, error) {
	product := new(models.Product)

	if err := tx.WithContext(ctx).
		Preload("Stocks", func(db *gorm.DB) *gorm.DB {
			if reqData.Filter.WarehouseID > 0 {
				db = db.Where("warehouse_id", reqData.Filter.WarehouseID).Limit(1)
			}
			return db.Select(`distinct on ("warehouse_id") warehouse_id,*`).
				Preload("Warehouse", func(db *gorm.DB) *gorm.DB {
					return db.Where("status", true)
				}).Order("warehouse_id,created_at desc")
		}).
		Where(&models.Product{
			ID: reqData.ID,
		}).
		Take(&product).Error; err != nil {
		r.Logger.Error(err)
		return nil, err
	}

	return product, nil
}

// FindAll :
func (r *ProductRepository) FindAll(ctx context.Context, reqData *models.Product, pagination paginate.Pagination, tx *gorm.DB) ([]*models.Product, *paginate.Pagination, error) {
	products := make([]*models.Product, 0)

	queryBuilder := tx.WithContext(ctx)

	if err := queryBuilder.Order("id desc").Scopes(paginate.Paginate(reqData, &pagination, queryBuilder)).
		Find(&products).Error; err != nil {
		r.Logger.Error(err)
		return nil, nil, err
	}

	return products, &pagination, nil
}
