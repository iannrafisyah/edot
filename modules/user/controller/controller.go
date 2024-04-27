package controller

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"Edot/models"
	"Edot/modules/user/dto"
	"Edot/modules/user/repository"
	"Edot/packages/logger"
	"Edot/utilities"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

type (
	IUserController interface {
		FindByID(ctx context.Context, reqData *dto.FindByIDRequest, tx *gorm.DB) (*dto.FindByIDResponse, error)
		FindByEmail(ctx context.Context, reqData *dto.FindByEmailRequest, tx *gorm.DB) (*dto.FindByEmailResponse, error)
	}

	UserController struct {
		fx.In
		UserRepository repository.IUserInterface
		Logger         *logger.Logger
	}
)

// NewController :
func NewController(userController UserController) IUserController {
	return &userController
}

// FindByEmail :
func (r *UserController) FindByEmail(ctx context.Context, reqData *dto.FindByEmailRequest, tx *gorm.DB) (*dto.FindByEmailResponse, error) {
	// Validate request data
	messages, err := reqData.Validate()
	if err != nil {
		r.Logger.Error(err)
		return nil, utilities.ErrorRequest(err, http.StatusBadRequest, messages)
	}

	fetchUser, err := r.UserRepository.FindByEmail(ctx, &models.User{
		Email: reqData.Email,
	}, tx)
	if err != nil {
		r.Logger.Error(err)
		if err == gorm.ErrRecordNotFound {
			return nil, utilities.ErrorRequest(
				fmt.Errorf(utilities.DataNotFound, "user"),
				http.StatusNotFound,
			)
		}
		return nil, utilities.ErrorRequest(
			errors.New(utilities.InternalServiceError),
			http.StatusInternalServerError,
		)
	}

	return &dto.FindByEmailResponse{
		ID:        fetchUser.ID,
		FullName:  fetchUser.FullName,
		Email:     fetchUser.Email,
		Password:  fetchUser.Password,
		CreatedAt: fetchUser.CreatedAt,
	}, nil
}

// FindByID :
func (r *UserController) FindByID(ctx context.Context, reqData *dto.FindByIDRequest, tx *gorm.DB) (*dto.FindByIDResponse, error) {
	fetchUser, err := r.UserRepository.FindByID(ctx, &models.User{
		ID: reqData.ContextUserID,
	}, tx)
	if err != nil {
		r.Logger.Error(err)
		if err == gorm.ErrRecordNotFound {
			return nil, utilities.ErrorRequest(
				fmt.Errorf(utilities.DataNotFound, "user"),
				http.StatusNotFound,
			)
		}
		return nil, utilities.ErrorRequest(
			errors.New(utilities.InternalServiceError),
			http.StatusInternalServerError,
		)
	}

	response := dto.FindByIDResponse{
		ID:        fetchUser.ID,
		FullName:  fetchUser.FullName,
		Email:     fetchUser.Email,
		CreatedAt: fetchUser.CreatedAt,
	}
	return &response, nil
}
