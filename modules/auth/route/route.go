package route

import (
	"net/http"

	"Edot/modules/auth/controller"
	auth "Edot/modules/auth/dto"
	"Edot/packages/logger"
	"Edot/packages/postgres"
	"Edot/routers"
	"Edot/utilities"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

type Handler struct {
	fx.In
	Controller controller.IAuthController
	Logger     *logger.Logger
	DB         *postgres.DB
	Router     *routers.Router
}

func NewRoute(h Handler, m ...echo.MiddlewareFunc) Handler {
	h.Route(m...)
	return h
}

func (r *Handler) Route(m ...echo.MiddlewareFunc) {
	echoRoute := r.Router.Group("/v1/auth", m...)
	echoRoute.POST("/login", r.Login)
}

// Login :
func (r *Handler) Login(c echo.Context) error {
	var reqData = new(auth.LoginRequest)

	if err := c.Bind(reqData); err != nil {
		r.Logger.Error(err)
		return utilities.Response(c, &utilities.ResponseRequest{
			StatusCode: http.StatusBadRequest,
			Message:    utilities.BadRequest,
		})
	}

	tx := r.DB.Gorm.Begin()
	resp, err := r.Controller.Login(c.Request().Context(), reqData, tx)
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
