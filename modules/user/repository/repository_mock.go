package repository

import (
	"Edot/models"
	"Edot/packages/logger"
	"context"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"
	"go.uber.org/fx"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type MockUserRepository struct {
	fx.In
	Logger *logger.Logger
	Mock   *mock.Mock
}

// NewMockRepository
func NewMockRepository(mockUserRepository MockUserRepository) IUserInterface {
	return &mockUserRepository
}

func (r *MockUserRepository) fakeUser() models.User {
	password, err := bcrypt.GenerateFromPassword([]byte("testing123"), bcrypt.DefaultCost)
	if err != nil {
		r.Logger.Error(err)
	}

	return models.User{
		ID:       1,
		FullName: gofakeit.Name(),
		Email:    "testing@mail.com",
		Password: string(password),
	}
}

// FindByID :
func (r *MockUserRepository) FindByID(ctx context.Context, reqData *models.User, tx *gorm.DB) (*models.User, error) {
	args := r.Mock.Called(reqData)
	request := args.Get(0).(models.User)
	fetchUser := r.fakeUser()

	if fetchUser.ID != request.ID {
		return nil, gorm.ErrRecordNotFound
	}

	return &fetchUser, nil
}

// FindByEmail :
func (r *MockUserRepository) FindByEmail(ctx context.Context, reqData *models.User, tx *gorm.DB) (*models.User, error) {
	args := r.Mock.Called(reqData)
	request := args.Get(0).(models.User)
	fetchUser := r.fakeUser()

	if fetchUser.Email != request.Email {
		return nil, gorm.ErrRecordNotFound
	}

	return &fetchUser, nil
}
