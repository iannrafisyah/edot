package route

import (
	"Edot/config"
	"Edot/models"
	"Edot/modules/auth/controller"
	userController "Edot/modules/user/controller"
	userRepository "Edot/modules/user/repository"
	"net/http"
	"net/http/httptest"
	"strings"

	"Edot/packages/logger"
	"Edot/packages/postgres"
	"Edot/routers"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestHandler_Login(t *testing.T) {
	type ModuleDependency struct {
		fx.In
		MockUserRepository userRepository.MockUserRepository
		AuthController     controller.IAuthController
		UserController     userController.IUserController
		Logger             *logger.Logger
		Handler            Handler
	}

	config.SetConfig("../../../")

	var moduleDependency ModuleDependency

	app := fxtest.New(t,
		fx.Provide(func() *mock.Mock {
			return &mock.Mock{}
		}),
		fx.Provide(postgres.NewMockPostgres),
		fx.Provide(routers.NewRouter),
		fx.Provide(logger.NewLogger),
		fx.Provide(userRepository.NewMockRepository),
		fx.Provide(controller.NewController),
		fx.Provide(userController.NewController),
		fx.Populate(&moduleDependency),
	)
	app.RequireStart().RequireStop()

	t.Run("LoginFailedForbidden", func(t *testing.T) {
		moduleDependency.MockUserRepository.Mock.On("FindByEmail", mock.Anything).Return(models.User{
			Email:    "testing@mail.com",
			Password: "testing124",
		}).Once()

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{
			"email":"testing@mail.com",
			"password":"testing124"
		}`))

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)

		ctx := c.Request().Context()
		c.SetRequest(c.Request().WithContext(ctx))

		if assert.NoError(t, moduleDependency.Handler.Login(c)) {
			assert.Equal(t, http.StatusForbidden, rec.Code)
		}
	})

	t.Run("LoginSuccess", func(t *testing.T) {
		moduleDependency.MockUserRepository.Mock.On("FindByEmail", mock.Anything).Return(models.User{
			Email:    "testing@mail.com",
			Password: "testing123",
		}).Once()

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{
			"email":"testing@mail.com",
			"password":"testing123"
		}`))

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)

		ctx := c.Request().Context()
		c.SetRequest(c.Request().WithContext(ctx))

		if assert.NoError(t, moduleDependency.Handler.Login(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

}
