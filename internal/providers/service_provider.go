package providers

import (
	hotelService "github.com/MaxFando/application-design/internal/core/hotel/service"
	orderService "github.com/MaxFando/application-design/internal/core/order/service"
)

type ServiceProvider struct {
	orderService *orderService.Service
	hotelService *hotelService.AvailabilityService
}

func NewServiceProvider() *ServiceProvider {
	return &ServiceProvider{}
}

func (sp *ServiceProvider) RegisterDependencies(repoProvider *RepositoryProvider) {
	sp.orderService = orderService.NewService(repoProvider.inlineOrderRepo)
	sp.hotelService = hotelService.NewAvailabilityService(repoProvider.inlineHotelRepo)
}
