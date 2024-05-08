package converter

import (
	"github.com/MaxFando/application-design/internal/core/order/entity"
)

func ToOrderCreateResponse(order entity.Order) OrderCreateResponse {
	return OrderCreateResponse{
		Order{
			HotelId:   order.HotelId,
			RoomId:    order.RoomIds[0],
			UserEmail: order.UserEmail,
			From:      order.From,
			To:        order.To,
		},
	}
}
