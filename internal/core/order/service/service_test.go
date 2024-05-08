package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/MaxFando/application-design/internal/core/order/entity"
	"github.com/MaxFando/application-design/internal/core/order/service"
)

func TestService_NewService(t *testing.T) {
	t.Run("should return new service", func(t *testing.T) {
		s := service.NewService(nil)
		assert.NotNil(t, s)
	})
}

func TestService_CreateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		ctx   context.Context
		order entity.Order
	}

	testCases := []struct {
		name         string
		initMockRepo func(m *service.MockRepository)
		err          error
		args         args
	}{
		{
			name: "success",
			initMockRepo: func(m *service.MockRepository) {
				m.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Return(nil)
			},
			args: args{
				ctx:   context.TODO(),
				order: entity.Order{},
			},
		},
		{
			name: "error",
			initMockRepo: func(m *service.MockRepository) {
				m.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Return(assert.AnError)
			},
			err: assert.AnError,
			args: args{
				ctx:   context.TODO(),
				order: entity.Order{},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := service.NewMockRepository(ctrl)
			tc.initMockRepo(mockRepo)
			s := service.NewService(mockRepo)
			err := s.CreateOrder(tc.args.ctx, tc.args.order)
			assert.Equal(t, tc.err, err)
		})
	}
}
