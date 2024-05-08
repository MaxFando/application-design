package usecase_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/MaxFando/application-design/internal/core/order/entity"
	"github.com/MaxFando/application-design/internal/core/order/usecase"
	"github.com/MaxFando/application-design/pkg/utils"
)

func TestMain(m *testing.M) {
	utils.InitializeLogger("info")

	os.Exit(m.Run())
}

func TestUseCase_NewUseCase(t *testing.T) {
	t.Run("should return new usecase", func(t *testing.T) {
		uc := usecase.NewUseCase(nil, nil)
		assert.NotNil(t, uc)
	})
}

func TestUseCase_CreateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		ctx   context.Context
		order entity.Order
	}

	testCases := []struct {
		name                 string
		args                 args
		initMockService      func(m *usecase.MockService)
		initMockHotelUseCase func(m *usecase.MockAvailabilityChecker)
		err                  error
	}{
		{
			name: "success",
			args: args{
				ctx:   context.TODO(),
				order: entity.Order{},
			},
			initMockService: func(m *usecase.MockService) {
				m.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Return(nil)
			},
			initMockHotelUseCase: func(m *usecase.MockAvailabilityChecker) {
				m.EXPECT().Check(gomock.Any(), gomock.Any()).Return(true, nil)
			},
		},
		{
			name: "error on check",
			args: args{
				ctx:   context.TODO(),
				order: entity.Order{},
			},
			initMockHotelUseCase: func(m *usecase.MockAvailabilityChecker) {
				m.EXPECT().Check(gomock.Any(), gomock.Any()).Return(false, assert.AnError)
			},
			err: assert.AnError,
		},
		{
			name: "error on create order",
			args: args{
				ctx:   context.TODO(),
				order: entity.Order{},
			},
			initMockService: func(m *usecase.MockService) {
				m.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Return(assert.AnError)
			},
			initMockHotelUseCase: func(m *usecase.MockAvailabilityChecker) {
				m.EXPECT().Check(gomock.Any(), gomock.Any()).Return(true, nil)
			},
			err: assert.AnError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService := usecase.NewMockService(ctrl)
			if tc.initMockService != nil {
				tc.initMockService(mockService)
			}

			mockHotelUseCase := usecase.NewMockAvailabilityChecker(ctrl)
			if tc.initMockHotelUseCase != nil {
				tc.initMockHotelUseCase(mockHotelUseCase)
			}

			uc := usecase.NewUseCase(mockService, mockHotelUseCase)
			err := uc.CreateOrder(tc.args.ctx, tc.args.order)
			assert.Equal(t, tc.err, err)
		})
	}
}
