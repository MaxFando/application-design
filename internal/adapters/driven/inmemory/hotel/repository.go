package hotel

import (
	"context"

	"github.com/pkg/errors"

	domainHotel "github.com/MaxFando/application-design/internal/core/hotel/entity"
	"github.com/MaxFando/application-design/internal/enum"
	"github.com/MaxFando/application-design/pkg/storage/inmemory"
)

var (
	ErrorNoDataInCache      = errors.New("hotel room: no data in cache")
	ErrorInvalidDataInCache = errors.New("hotel room: invalid data in cache")
)

type Repository struct {
	mem inmemory.Cache
}

func NewRepository(mem inmemory.Cache) *Repository {
	return &Repository{
		mem: mem,
	}
}

func (r *Repository) InitializeAvailability(ctx context.Context, availability []domainHotel.RoomAvailability) error {
	return r.mem.Set(ctx, enum.InlineAvailabilityCacheKey, availability)
}

func (r *Repository) GetAvailability(ctx context.Context) ([]domainHotel.RoomAvailability, error) {
	mem, ok := r.mem.Get(ctx, enum.InlineAvailabilityCacheKey)
	if !ok {
		return nil, ErrorNoDataInCache
	}

	availability, ok := mem.([]domainHotel.RoomAvailability)
	if !ok {
		return nil, ErrorInvalidDataInCache
	}

	return availability, nil
}

func (r *Repository) UpdateAvailability(ctx context.Context, idx int, availability domainHotel.RoomAvailability) error {
	availabilityList, err := r.GetAvailability(ctx)
	if err != nil {
		return err
	}

	availabilityList[idx] = availability

	return r.mem.Set(ctx, enum.InlineAvailabilityCacheKey, availabilityList)
}
