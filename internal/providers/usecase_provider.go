package providers

import (
	hotelUseCase "github.com/MaxFando/application-design/internal/core/hotel/usecase"
	orderuseCase "github.com/MaxFando/application-design/internal/core/order/usecase"
)

type UseCaseProvider struct {
	HotelUseCase *hotelUseCase.AvailabilityUseCase
	OrderUseCase *orderuseCase.UseCase
}

func NewUseCaseProvider() *UseCaseProvider {
	return &UseCaseProvider{}
}

func (ucp *UseCaseProvider) RegisterDependencies(serviceProvider *ServiceProvider) {
	ucp.HotelUseCase = hotelUseCase.NewAvailabilityUseCase(serviceProvider.hotelService)
	ucp.OrderUseCase = orderuseCase.NewUseCase(serviceProvider.orderService, ucp.HotelUseCase)
}
