package service

import (
	"context"

	domainHotel "github.com/MaxFando/application-design/internal/core/hotel/entity"
)

//go:generate mockgen -source=$GOFILE -destination=./mock_service_ports.go -package=${GOPACKAGE}

type AvailabilityInitializerFetcherRepository interface {
}

type AvailabilityInitializerWriterRepository interface {
	InitializeAvailability(ctx context.Context, availability []domainHotel.RoomAvailability) error
}

type AvailabilityRepository interface {
	GetAvailability(ctx context.Context) ([]domainHotel.RoomAvailability, error)
	UpdateAvailability(ctx context.Context, idx int, availability domainHotel.RoomAvailability) error
}
