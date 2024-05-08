package service

import (
	"context"

	"github.com/MaxFando/application-design/internal/core/order/entity"
)

//go:generate mockgen -source=$GOFILE -destination=./mock_ports.go -package=${GOPACKAGE}

type Repository interface {
	CreateOrder(ctx context.Context, order entity.Order) error
}
