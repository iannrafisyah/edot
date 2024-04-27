package route

import (
	"net/http"

	"Edot/modules/user/controller"
	"Edot/modules/user/dto"
	"Edot/packages/jwt"
	"Edot/packages/logger"
	"Edot/packages/postgres"
	"Edot/routers"
	"Edot/utilities"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

type Handler struct {
	fx.In
	Controller controller.IUserController
	Logger     *logger.Logger
	DB         *postgres.DB
	Router     *routers.Router
}

func NewRoute(h Handler, m ...echo.MiddlewareFunc) Handler {
	h.Route(m...)
	return h
}

func (r *Handler) Route(m ...echo.MiddlewareFunc) {
	echoRoute := r.Router.Group("/v1/user", m...)
	echoRoute.Use(r.Router.AuthGuard)
	echoRoute.GET("/profile", r.Detail)
}

// Detail :
func (r *Handler) Detail(c echo.Context) error {
	var reqData = new(dto.FindByIDRequest)

	data, ok := c.Request().Context().Value(jwt.InternalClaimData{}).(jwt.InternalClaimData)
	if !ok {
		return utilities.Response(c, &utilities.ResponseRequest{
			StatusCode: http.StatusUnauthorized,
			Message:    utilities.Authorization,
		})
	}

	reqData.ContextUserID = data.UserID

	tx := r.DB.Gorm.Begin()
	resp, err := r.Controller.FindByID(c.Request().Context(), reqData, tx)
	if err != nil {
		r.Logger.Error(err)
		defer func() {
			tx.Rollback()
		}()

		return utilities.Response(c, &utilities.ResponseRequest{
			StatusCode: utilities.ParseError(err).StatusCode,
			Data:       utilities.ParseError(err).Data,
			Message:    err.Error(),
		})
	}
	tx.Commit()

	return utilities.Response(c, &utilities.ResponseRequest{
		StatusCode: http.StatusOK,
		Message:    utilities.Success,
		Data:       resp,
	})
}
