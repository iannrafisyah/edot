package route

import (
	"net/http"

	"Edot/modules/product/controller"
	"Edot/modules/product/dto"
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
	Controller controller.IProductController
	Logger     *logger.Logger
	DB         *postgres.DB
	Router     *routers.Router
}

func NewRoute(h Handler, m ...echo.MiddlewareFunc) Handler {
	h.Route(m...)
	return h
}

func (r *Handler) Route(m ...echo.MiddlewareFunc) {
	echoRoute := r.Router.Group("/v1/product", m...)
	echoRoute.Use(r.Router.AuthGuard)
	echoRoute.GET("", r.FindAll)
	echoRoute.GET("/:id", r.FindByID)
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

	if err := echo.QueryParamsBinder(c).
		Int("limit", &reqData.Limit).
		Int("page", &reqData.Page).
		BindError(); err != nil {
		r.Logger.Error(err)
		return utilities.Response(c, &utilities.ResponseRequest{
			StatusCode: http.StatusBadRequest,
			Message:    utilities.BadRequest,
		})
	}

	tx := r.DB.Gorm.Begin()

	resp, paging, err := r.Controller.FindAll(c.Request().Context(), reqData, tx)
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

	if paging != nil {
		paging.Next(c)
		paging.Prev(c)
	}

	return utilities.Response(c, &utilities.ResponseRequest{
		StatusCode: http.StatusOK,
		Message:    utilities.Success,
		Data:       resp,
		Paginate:   paging,
	})
}

// FindByID :
func (r *Handler) FindByID(c echo.Context) error {
	var reqData = new(dto.FindByIDRequest)

	data, ok := c.Request().Context().Value(jwt.InternalClaimData{}).(jwt.InternalClaimData)
	if !ok {
		return utilities.Response(c, &utilities.ResponseRequest{
			StatusCode: http.StatusUnauthorized,
			Message:    utilities.Authorization,
		})
	}

	reqData.ContextUserID = data.UserID

	if err := echo.PathParamsBinder(c).
		Int("id", &reqData.ID).
		BindError(); err != nil {
		r.Logger.Error(err)
		return utilities.Response(c, &utilities.ResponseRequest{
			StatusCode: http.StatusBadRequest,
			Message:    utilities.BadRequest,
		})
	}

	if err := echo.QueryParamsBinder(c).
		Int("warehouse_id", &reqData.WarehouseID).
		BindError(); err != nil {
		r.Logger.Error(err)
		return utilities.Response(c, &utilities.ResponseRequest{
			StatusCode: http.StatusBadRequest,
			Message:    utilities.BadRequest,
		})
	}

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
