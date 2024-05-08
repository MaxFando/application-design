package http

import (
	"context"

	"github.com/labstack/echo/v4"

	orderV1 "github.com/MaxFando/application-design/internal/adapters/driving/http/v1/order"
	"github.com/MaxFando/application-design/internal/providers"
)

func NewRouter(ctx context.Context, handler *echo.Echo) *echo.Echo {
	useCaseProvider := ctx.Value(providers.UseCaseProviderKey).(*providers.UseCaseProvider)

	orderController := orderV1.NewController(useCaseProvider.OrderUseCase)
	handler.POST("/orders", orderController.CreateOrder)

	return handler
}
