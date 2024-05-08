package entity

import (
	"strings"
	"time"
)

type UnavailableDays map[time.Time]struct{}

func (ud *UnavailableDays) FormatHumanReadable() string {
	days := make([]string, 0)
	for day := range *ud {
		days = append(days, day.Format("2006-01-02"))
	}

	return strings.Join(days, ", ")
}

type RoomAvailability struct {
	HotelID string    `json:"hotel_id"`
	RoomID  string    `json:"room_id"`
	Date    time.Time `json:"date"`
	Quota   int       `json:"quota"`
}
