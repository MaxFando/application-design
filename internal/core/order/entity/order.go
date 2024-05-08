package entity

import "time"

type Order struct {
	HotelId   string    `json:"hotel_id"`
	RoomIds   []string  `json:"room_id"`
	UserEmail string    `json:"email"`
	From      time.Time `json:"from"`
	To        time.Time `json:"to"`
}
