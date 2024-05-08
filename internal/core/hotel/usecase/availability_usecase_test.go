package usecase_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/MaxFando/application-design/internal/core/hotel/entity"
	"github.com/MaxFando/application-design/internal/core/hotel/usecase"
	"github.com/MaxFando/application-design/internal/tools"
)

func TestAvailabilityUseCase_NewUseCase(t *testing.T) {
	t.Run("Test NewUseCase", func(t *testing.T) {
		u := usecase.NewAvailabilityUseCase(nil)
		assert.NotNil(t, u)
	})
}

func TestAvailabilityUseCase_Check(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		daysToBook []time.Time
	}

	today := tools.Date(2021, 1, 1)
	tomorrow := tools.Date(2021, 1, 2)

	testCases := []struct {
		name            string
		args            args
		initMockService func(m *usecase.MockAvailabilityServiceInterface)
		want            bool
		err             error
	}{
		{
			name: "success",
			args: args{
				daysToBook: []time.Time{today, tomorrow},
			},
			initMockService: func(m *usecase.MockAvailabilityServiceInterface) {
				m.EXPECT().ComputeUnavailableDays(gomock.Any(), gomock.Any()).Return(nil, nil)
			},
			want: true,
		},
		{
			name: "error on service",
			args: args{
				daysToBook: []time.Time{today, tomorrow},
			},
			initMockService: func(m *usecase.MockAvailabilityServiceInterface) {
				m.EXPECT().ComputeUnavailableDays(gomock.Any(), gomock.Any()).Return(nil, assert.AnError)
			},
			want: false,
			err:  assert.AnError,
		},
		{
			name: "room not available",
			args: args{
				daysToBook: []time.Time{today, tomorrow},
			},
			initMockService: func(m *usecase.MockAvailabilityServiceInterface) {
				m.EXPECT().ComputeUnavailableDays(gomock.Any(), gomock.Any()).Return(entity.UnavailableDays{
					today: struct{}{},
				}, nil)
			},
			want: false,
			err:  usecase.ErrRoomNotAvailable,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			service := usecase.NewMockAvailabilityServiceInterface(ctrl)
			tt.initMockService(service)

			uc := usecase.NewAvailabilityUseCase(service)
			got, err := uc.Check(nil, tt.args.daysToBook)

			assert.Equal(t, tt.want, got)
			if err != nil {
				assert.ErrorAs(t, err, &tt.err)
			}
		})
	}
}
