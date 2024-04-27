package route

import (
	"net/http"

	"Edot/modules/cart/controller"
	"Edot/modules/cart/dto"
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
	Controller controller.ICartController
	Logger     *logger.Logger
	DB         *postgres.DB
	Router     *routers.Router
}

func NewRoute(h Handler, m ...echo.MiddlewareFunc) Handler {
	h.Route(m...)
	return h
}

func (r *Handler) Route(m ...echo.MiddlewareFunc) {
	echoRoute := r.Router.Group("/v1/cart", m...)
	echoRoute.Use(r.Router.AuthGuard)
	echoRoute.GET("", r.FindAll)
	echoRoute.POST("", r.AddCart)
}

// FindAll :
func (r *Handler) FindAll(c echo.Context) error {
	var reqData = new(dto.FindAllRequest)

	data, ok := c.Request().Context().Value(jwt.InternalClaimData{}).(jwt.InternalClaimData)
	if !ok {
		return utilities.Response(c, &utilities.ResponseRequest{
			StatusCode: http.StatusUnauthorized,
			Message:    utilities.Authorization,
		})
	}

	reqData.ContextUserID = data.UserID

	tx := r.DB.Gorm.Begin()
	resp, err := r.Controller.FindAll(c.Request().Context(), reqData, tx)
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

// AddCart :
func (r *Handler) AddCart(c echo.Context) error {
	var reqData = new(dto.AddCartRequest)

	data, ok := c.Request().Context().Value(jwt.InternalClaimData{}).(jwt.InternalClaimData)
	if !ok {
		return utilities.Response(c, &utilities.ResponseRequest{
			StatusCode: http.StatusUnauthorized,
			Message:    utilities.Authorization,
		})
	}

	if err := c.Bind(reqData); err != nil {
		r.Logger.Error(err)
		return utilities.Response(c, &utilities.ResponseRequest{
			StatusCode: http.StatusBadRequest,
			Message:    utilities.BadRequest,
		})
	}

	reqData.ContextUserID = data.UserID

	tx := r.DB.Gorm.Begin()
	if err := r.Controller.AddCart(c.Request().Context(), reqData, tx); err != nil {
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
		StatusCode: http.StatusCreated,
		Message:    utilities.Success,
	})
}
