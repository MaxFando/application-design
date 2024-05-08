package order

import (
	"context"
	"sync"

	"github.com/pkg/errors"

	"github.com/MaxFando/application-design/internal/core/order/entity"
	"github.com/MaxFando/application-design/internal/enum"
	"github.com/MaxFando/application-design/pkg/storage/inmemory"
)

var (
	ErrorInvalidDataInCache = errors.New("order: invalid data in cache")
)

type Repository struct {
	mem inmemory.Cache
	mu  sync.Mutex
}

func NewRepository(mem inmemory.Cache) *Repository {
	return &Repository{
		mem: mem,
	}
}

func (r *Repository) GetOrders(ctx context.Context) ([]entity.Order, error) {
	data, exists := r.mem.Get(ctx, enum.InlineOrdersCacheKey)
	if !exists {
		return make([]entity.Order, 0), nil
	}

	orders, ok := data.([]entity.Order)
	if !ok {
		return nil, ErrorInvalidDataInCache
	}

	return orders, nil
}

func (r *Repository) CreateOrder(ctx context.Context, order entity.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	orders, err := r.GetOrders(ctx)
	if err != nil {
		return err
	}

	orders = append(orders, order)
	return r.mem.Set(ctx, enum.InlineOrdersCacheKey, orders)
}
