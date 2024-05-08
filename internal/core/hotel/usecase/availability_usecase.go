package usecase

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/MaxFando/application-design/internal/core/hotel/entity"
)

//go:generate mockgen -source=$GOFILE -destination=./mock_service.go -package=${GOPACKAGE}

var ErrRoomNotAvailable = errors.New("Hotel room is not available for selected dates")

type AvailabilityServiceInterface interface {
	ComputeUnavailableDays(ctx context.Context, daysToBook []time.Time) (entity.UnavailableDays, error)
}

type AvailabilityUseCase struct {
	service AvailabilityServiceInterface
}

func NewAvailabilityUseCase(service AvailabilityServiceInterface) *AvailabilityUseCase {
	return &AvailabilityUseCase{service: service}
}

func (uc *AvailabilityUseCase) Check(ctx context.Context, daysToBook []time.Time) (bool, error) {
	unavailableDays, err := uc.service.ComputeUnavailableDays(ctx, daysToBook)
	if err != nil {
		return false, err
	}

	if len(unavailableDays) != 0 {
		return false, errors.Wrapf(ErrRoomNotAvailable, "unavailable days: %s", unavailableDays.FormatHumanReadable())
	}

	return true, nil
}
