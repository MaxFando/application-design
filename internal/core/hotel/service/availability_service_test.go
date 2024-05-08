package service_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/MaxFando/application-design/internal/core/hotel/entity"
	"github.com/MaxFando/application-design/internal/core/hotel/service"
	"github.com/MaxFando/application-design/internal/tools"
	"github.com/MaxFando/application-design/pkg/utils"
)

func TestMain(m *testing.M) {
	utils.InitializeLogger("info")

	os.Exit(m.Run())
}

func TestAvailabilityService_NewAvailabilityService(t *testing.T) {
	t.Run("should return new availability service", func(t *testing.T) {
		s := service.NewAvailabilityService(nil)
		assert.NotNil(t, s)
	})
}

func TestAvailabilityService_ComputeUnavailableDays(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		ctx        context.Context
		daysToBook []time.Time
	}

	today := tools.Date(2021, 1, 1)
	tomorrow := tools.Date(2021, 1, 2)

	testCases := []struct {
		name               string
		args               args
		initMockRepository func(m *service.MockAvailabilityRepository)
		err                error
		want               entity.UnavailableDays
	}{
		{
			name: "success",
			args: args{
				ctx:        context.TODO(),
				daysToBook: []time.Time{today, tomorrow},
			},
			initMockRepository: func(m *service.MockAvailabilityRepository) {
				roomToday := entity.RoomAvailability{Date: today, Quota: 1}
				roomTomorrow := entity.RoomAvailability{Date: tomorrow, Quota: 1}

				m.EXPECT().GetAvailability(gomock.Any()).Return([]entity.RoomAvailability{roomToday, roomTomorrow}, nil)

				m.EXPECT().UpdateAvailability(gomock.Any(), 0, entity.RoomAvailability{Date: today, Quota: 0}).Return(nil)
				m.EXPECT().UpdateAvailability(gomock.Any(), 1, entity.RoomAvailability{Date: tomorrow, Quota: 0}).Return(nil)
			},
			want: entity.UnavailableDays{},
		},
		{
			name: "error on get availability",
			args: args{
				ctx:        context.TODO(),
				daysToBook: []time.Time{today, tomorrow},
			},
			initMockRepository: func(m *service.MockAvailabilityRepository) {
				m.EXPECT().GetAvailability(gomock.Any()).Return(nil, assert.AnError)
			},
			err: assert.AnError,
		},
		{
			name: "error on update availability",
			args: args{
				ctx:        context.TODO(),
				daysToBook: []time.Time{today, tomorrow},
			},
			initMockRepository: func(m *service.MockAvailabilityRepository) {
				roomToday := entity.RoomAvailability{Date: today, Quota: 1}
				roomTomorrow := entity.RoomAvailability{Date: tomorrow, Quota: 1}

				m.EXPECT().GetAvailability(gomock.Any()).Return([]entity.RoomAvailability{roomToday, roomTomorrow}, nil)

				m.EXPECT().UpdateAvailability(gomock.Any(), 0, entity.RoomAvailability{Date: today, Quota: 0}).Return(nil)
				m.EXPECT().UpdateAvailability(gomock.Any(), 1, entity.RoomAvailability{Date: tomorrow, Quota: 0}).Return(assert.AnError)
			},
			err: assert.AnError,
		},
		{
			name: "success with unavailable days",
			args: args{
				ctx:        context.TODO(),
				daysToBook: []time.Time{today, tomorrow},
			},
			initMockRepository: func(m *service.MockAvailabilityRepository) {
				roomToday := entity.RoomAvailability{Date: today, Quota: 0}
				roomTomorrow := entity.RoomAvailability{Date: tomorrow, Quota: 1}

				m.EXPECT().GetAvailability(gomock.Any()).Return([]entity.RoomAvailability{roomToday, roomTomorrow}, nil)

				m.EXPECT().UpdateAvailability(gomock.Any(), 1, entity.RoomAvailability{Date: tomorrow, Quota: 0}).Return(nil)
			},
			want: entity.UnavailableDays{today: struct{}{}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepository := service.NewMockAvailabilityRepository(ctrl)
			if tc.initMockRepository != nil {
				tc.initMockRepository(mockRepository)
			}

			s := service.NewAvailabilityService(mockRepository)
			got, err := s.ComputeUnavailableDays(tc.args.ctx, tc.args.daysToBook)
			assert.Equal(t, tc.err, err)
			assert.Equal(t, len(tc.want), len(got))
		})
	}
}
