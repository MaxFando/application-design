package http

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	middlewares "github.com/MaxFando/application-design/internal/adapters/driving/http/middleware"
)

type customValidator struct {
	Validator *validator.Validate
}

func (cv *customValidator) Validate(i interface{}) error {
	err := cv.Validator.Struct(i)

	return err
}

type Handler struct {
	echo *echo.Echo
}

func NewHandler(ctx context.Context) *echo.Echo {
	echoMainServer := echo.New()
	echoMainServer.Validator = &customValidator{Validator: validator.New()}
	echoMainServer.Use(curryContextMiddleware(ctx))
	echoMainServer.Use(middlewares.Logger)

	echoMainServer = NewRouter(ctx, echoMainServer)

	return echoMainServer
}

func curryContextMiddleware(ctx context.Context) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}
