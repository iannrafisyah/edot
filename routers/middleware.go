package routers

import (
	"context"
	"net/http"
	"strings"

	"Edot/config"
	"Edot/models"
	"Edot/packages/jwt"
	"Edot/utilities"

	"github.com/labstack/echo/v4"
)

// AuthGuard :
func (r *Router) AuthGuard(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authorizationHeader := c.Request().Header.Get("Authorization")
		authorization := strings.Split(authorizationHeader, " ")
		if len(authorization) > 1 {
			result, err := jwt.ParseClaim(authorization[1], config.Get().Auth.Secret)

			if err != nil {
				r.Logger.Error(err)
				return utilities.Response(c, &utilities.ResponseRequest{
					StatusCode: http.StatusUnauthorized,
					Message:    utilities.Authorization,
				})
			}

			var ctx = c.Request().Context()

			// Check exist user
			tx := r.db.Gorm.Begin()
			fetchUser, err := r.userRepository.FindByID(ctx, &models.User{
				ID: result.Data.UserID,
			}, tx)
			if err != nil {
				r.Logger.Error(err)
				defer func() {
					tx.Rollback()
				}()

				return utilities.Response(c, &utilities.ResponseRequest{
					StatusCode: http.StatusUnauthorized,
					Message:    utilities.Authorization,
				})
			}
			tx.Commit()

			jwtClaimData := jwt.InternalClaimData{
				UserID: fetchUser.ID,
			}

			ctx = context.WithValue(ctx, jwt.InternalClaimData{}, jwtClaimData)

			c.SetRequest(c.Request().WithContext(ctx))
		} else {
			return utilities.Response(c, &utilities.ResponseRequest{
				StatusCode: http.StatusUnauthorized,
				Message:    utilities.Authorization,
			})
		}
		return next(c)
	}
}
