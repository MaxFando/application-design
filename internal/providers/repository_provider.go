package providers

import (
	"context"

	inmemoryHotel "github.com/MaxFando/application-design/internal/adapters/driven/inmemory/hotel"
	inmemoryOrder "github.com/MaxFando/application-design/internal/adapters/driven/inmemory/order"
	domainHotel "github.com/MaxFando/application-design/internal/core/hotel/entity"
	"github.com/MaxFando/application-design/internal/tools"
	"github.com/MaxFando/application-design/pkg/storage/inmemory"
)

type RepositoryProvider struct {
	inlineOrderRepo *inmemoryOrder.Repository
	inlineHotelRepo *inmemoryHotel.Repository
}

func NewRepositoryProvider() *RepositoryProvider {
	return &RepositoryProvider{}
}

func (rp *RepositoryProvider) RegisterDependencies() {
	rp.inlineOrderRepo = inmemoryOrder.NewRepository(inmemory.New())
	rp.inlineHotelRepo = inmemoryHotel.NewRepository(inmemory.New())

	var Availability = []domainHotel.RoomAvailability{
		{"reddison", "lux", tools.Date(2024, 1, 1), 1},
		{"reddison", "lux", tools.Date(2024, 1, 2), 1},
		{"reddison", "lux", tools.Date(2024, 1, 3), 1},
		{"reddison", "lux", tools.Date(2024, 1, 4), 1},
		{"reddison", "lux", tools.Date(2024, 1, 5), 0},
	}

	// Костыльное решение, но в рамках тестового задания сойдет
	_ = rp.inlineHotelRepo.InitializeAvailability(context.TODO(), Availability)
}
