package usecase

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	hotelUseCase "github.com/MaxFando/application-design/internal/core/hotel/usecase"
	"github.com/MaxFando/application-design/internal/core/order/entity"
	"github.com/MaxFando/application-design/internal/tools"
	"github.com/MaxFando/application-design/pkg/utils"
)

//go:generate mockgen -source=$GOFILE -destination=./mock_service.go -package=${GOPACKAGE}

type Service interface {
	CreateOrder(ctx context.Context, order entity.Order) error
}

type AvailabilityChecker interface {
	Check(ctx context.Context, daysToBook []time.Time) (bool, error)
}

type UseCase struct {
	service      Service
	hotelUseCase AvailabilityChecker

	// сюда можно добавить отдельный usecase для отправка письма-подтверждения о бронировании
	// скидки, промокоды, программы лояльности
}

func NewUseCase(service Service, hotelUseCase AvailabilityChecker) *UseCase {
	return &UseCase{
		service:      service,
		hotelUseCase: hotelUseCase,
	}
}

func (uc *UseCase) CreateOrder(ctx context.Context, newOrder entity.Order) error {
	daysToBook := tools.DaysBetween(newOrder.From, newOrder.To)
	_, err := uc.hotelUseCase.Check(ctx, daysToBook)
	if err != nil {
		if errors.Is(err, hotelUseCase.ErrRoomNotAvailable) {
			utils.Logger.ErrorWithContext(ctx)(hotelUseCase.ErrRoomNotAvailable, zap.Any("order", newOrder))
		}

		return err
	}

	err = uc.service.CreateOrder(ctx, newOrder)
	if err != nil {
		return err
	}

	utils.Logger.InfoWithContext(ctx)("Order successfully created", zap.Any("order", newOrder))
	return nil
}
