package controller

import (
	"context"
	"errors"
	"net/http"
	"time"

	"Edot/config"
	"Edot/modules/auth/dto"
	userController "Edot/modules/user/controller"
	userDto "Edot/modules/user/dto"
	_jwt "Edot/packages/jwt"
	"Edot/packages/logger"
	"Edot/utilities"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt"
	"go.uber.org/fx"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type (
	IAuthController interface {
		Login(ctx context.Context, reqData *dto.LoginRequest, tx *gorm.DB) (*dto.LoginResponse, error)
	}

	AuthController struct {
		fx.In
		Logger         *logger.Logger
		UserController userController.IUserController
	}
)

// NewController :
func NewController(authController AuthController) IAuthController {
	return &authController
}

// Login :
func (r *AuthController) Login(ctx context.Context, reqData *dto.LoginRequest, tx *gorm.DB) (*dto.LoginResponse, error) {
	// Validate request data
	messages, err := reqData.Validate()
	if err != nil {
		r.Logger.Error(err)
		return nil, utilities.ErrorRequest(err, http.StatusBadRequest, messages)
	}

	fetchUser, err := r.UserController.FindByEmail(ctx, &userDto.FindByEmailRequest{
		Email: reqData.Email,
	}, tx)
	if err != nil {
		r.Logger.Error(err)

		if utilities.ParseError(err).StatusCode == http.StatusNotFound {
			return nil, utilities.ErrorRequest(errors.New(utilities.InvalidAccessLogin), http.StatusForbidden)
		}

		return nil, err
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(fetchUser.Password), []byte(reqData.Password)); err != nil {
		r.Logger.Error(err)
		return nil, utilities.ErrorRequest(errors.New(utilities.InvalidAccessLogin), http.StatusForbidden)
	}

	// Generate uuid for user jwt
	generateUUID, err := uuid.NewV4()
	if err != nil {
		r.Logger.Error(err)
		return nil, utilities.ErrorRequest(err, http.StatusInternalServerError)
	}

	// Generate access and refresh token
	accessToken, err := _jwt.GenerateToken(_jwt.Claim{
		Data: _jwt.ClaimData{
			UserID: fetchUser.ID,
			UUID:   generateUUID.String(),
		},
		StandardClaims: jwt.StandardClaims{
			Audience:  "", // Web | Mobile = Get from context header
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(config.Get().Auth.ExpireAccessTokenDuration).Unix(),
		},
	}, config.Get().Auth.Secret)
	if err != nil {
		r.Logger.Error(err)
		return nil, utilities.ErrorRequest(err, http.StatusInternalServerError)
	}

	refreshToken, err := _jwt.GenerateToken(_jwt.Claim{
		Data: _jwt.ClaimData{
			UserID: fetchUser.ID,
			UUID:   generateUUID.String(),
		},
		StandardClaims: jwt.StandardClaims{
			Audience:  "", // Web | Mobile = Get from context header
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(config.Get().Auth.ExpireRefreshTokenDuration).Unix(),
		},
	}, config.Get().Auth.SecretClaim)
	if err != nil {
		r.Logger.Error(err)
		return nil, utilities.ErrorRequest(err, http.StatusInternalServerError)
	}

	return &dto.LoginResponse{
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
	}, nil
}
