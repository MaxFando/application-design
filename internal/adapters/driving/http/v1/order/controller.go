package order

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/MaxFando/application-design/internal/adapters/driving/http/v1/order/converter"
	"github.com/MaxFando/application-design/internal/core/order/entity"
	orderUseCase "github.com/MaxFando/application-design/internal/core/order/usecase"
	"github.com/MaxFando/application-design/internal/tools"
	"github.com/MaxFando/application-design/pkg/utils"
)

type Controller struct {
	uc *orderUseCase.UseCase
}

func NewController(uc *orderUseCase.UseCase) *Controller {
	return &Controller{
		uc: uc,
	}
}

type requestCreateOrder struct {
	HotelId string    `json:"hotel_id" validate:"required"`
	RoomId  string    `json:"room_id" validate:"required"`
	Email   string    `json:"email" validate:"required"`
	From    time.Time `json:"from" validate:"required"`
	To      time.Time `json:"to" validate:"required"`
}

func (req *requestCreateOrder) toDomain() entity.Order {
	return entity.Order{
		HotelId:   req.HotelId,
		RoomIds:   []string{req.RoomId},
		UserEmail: req.Email,
		From:      req.From,
		To:        req.To,
	}
}

func (ctr *Controller) CreateOrder(c echo.Context) error {
	ctx := c.Request().Context()

	request := new(requestCreateOrder)
	if err := c.Bind(request); err != nil {
		return tools.SendEchoHttpError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(request); err != nil {
		return tools.SendEchoHttpError(http.StatusUnprocessableEntity, err.Error())
	}

	newOrder := request.toDomain()
	err := ctr.uc.CreateOrder(ctx, newOrder)
	if err != nil {
		return tools.SendEchoHttpError(http.StatusConflict, err.Error())
	}

	utils.Logger.InfoWithContext(ctx)("Order successfully created", zap.Any("order", newOrder))
	return c.JSON(http.StatusCreated, converter.ToOrderCreateResponse(newOrder))
}
