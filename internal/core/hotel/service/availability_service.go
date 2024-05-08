package service

import (
	"context"
	"sync"
	"time"

	"github.com/MaxFando/application-design/internal/core/hotel/entity"
)

type AvailabilityService struct {
	repository AvailabilityRepository
	mu         sync.Mutex
}

func NewAvailabilityService(repository AvailabilityRepository) *AvailabilityService {
	return &AvailabilityService{
		repository: repository,
	}
}

func (s *AvailabilityService) ComputeUnavailableDays(ctx context.Context, daysToBook []time.Time) (entity.UnavailableDays, error) {
	// todo
	// тут придется все обернуть в mutex так как доступ к данным из разных горутин
	// и разлочить придется в самом конце, возможно позже придумаю способ получше
	s.mu.Lock()
	defer s.mu.Unlock()

	unavailableDays := make(map[time.Time]struct{})
	for _, day := range daysToBook {
		unavailableDays[day] = struct{}{}
	}

	availabilities, err := s.repository.GetAvailability(ctx)
	if err != nil {
		return nil, err
	}

	for _, dayToBook := range daysToBook {
		for idx, availability := range availabilities {
			if !availability.Date.Equal(dayToBook) || availability.Quota < 1 {
				continue
			}

			availability.Quota -= 1
			errStore := s.repository.UpdateAvailability(ctx, idx, availability)
			if errStore != nil {
				return nil, errStore
			}
			delete(unavailableDays, dayToBook)
		}
	}

	return unavailableDays, nil
}
