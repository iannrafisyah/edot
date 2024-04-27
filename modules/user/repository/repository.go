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
	IUserInterface interface {
		FindByID(ctx context.Context, reqData *models.User, tx *gorm.DB) (*models.User, error)
		FindByEmail(ctx context.Context, reqData *models.User, tx *gorm.DB) (*models.User, error)
	}

	UserRepository struct {
		fx.In
		DB     *postgres.DB
		Logger *logger.Logger
	}
)

// NewRepository :
func NewRepository(userRepository UserRepository) IUserInterface {
	return &userRepository
}

// FindByID :
func (r *UserRepository) FindByID(ctx context.Context, reqData *models.User, tx *gorm.DB) (*models.User, error) {
	user := new(models.User)

	if err := tx.WithContext(ctx).
		Where(&models.User{
			ID: reqData.ID,
		}).
		Take(&user).Error; err != nil {
		r.Logger.Error(err)
		return nil, err
	}

	return user, nil
}

// FindByEmail :
func (r *UserRepository) FindByEmail(ctx context.Context, reqData *models.User, tx *gorm.DB) (*models.User, error) {
	user := new(models.User)

	if err := tx.WithContext(ctx).
		Where(&models.User{
			Email: reqData.Email,
		}).
		Take(&user).Error; err != nil {
		r.Logger.Error(err)
		return nil, err
	}

	return user, nil
}
