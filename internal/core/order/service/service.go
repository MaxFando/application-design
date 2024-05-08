package service

import (
	"context"

	"github.com/MaxFando/application-design/internal/core/order/entity"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) CreateOrder(ctx context.Context, order entity.Order) error {
	return s.repository.CreateOrder(ctx, order)
}
