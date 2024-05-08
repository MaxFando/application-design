package converter

import "time"

type OrderCreateResponse struct {
	Order
}

type Order struct {
	HotelId   string    `json:"hotel_id"`
	RoomId    string    `json:"room_id"`
	UserEmail string    `json:"email"`
	From      time.Time `json:"from"`
	To        time.Time `json:"to"`
}
