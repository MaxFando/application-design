// Ниже реализован сервис бронирования номеров в отеле. В предметной области
// выделены два понятия: Order — заказ, который включает в себя даты бронирования
// и контакты пользователя, и RoomAvailability — количество свободных номеров на
// конкретный день.
//
// Задание:
// - провести рефакторинг кода с выделением слоев и абстракций
// - применить best-practices там где это имеет смысл
// - исправить имеющиеся в реализации логические и технические ошибки и неточности
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"

	domainHotel "github.com/MaxFando/application-design/internal/core/hotel/entity"
	domainOrder "github.com/MaxFando/application-design/internal/core/order/entity"
	"github.com/MaxFando/application-design/internal/tools"
	"github.com/MaxFando/application-design/pkg/utils"
)

var ErrRoomNotAvailable = errors.New("Hotel room is not available for selected dates")
var Orders = []domainOrder.Order{}

var Availability = []domainHotel.RoomAvailability{
	{"reddison", "lux", tools.Date(2024, 1, 1), 1},
	{"reddison", "lux", tools.Date(2024, 1, 2), 1},
	{"reddison", "lux", tools.Date(2024, 1, 3), 1},
	{"reddison", "lux", tools.Date(2024, 1, 4), 1},
	{"reddison", "lux", tools.Date(2024, 1, 5), 0},
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/orders", createOrder)

	LogInfo("Server listening on localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	if errors.Is(err, http.ErrServerClosed) {
		LogInfo("Server closed")
	} else if err != nil {
		LogErrorf("Server failed: %s", err)
		os.Exit(1)
	}
}

func createOrder(w http.ResponseWriter, r *http.Request) {
	var newOrder domainOrder.Order
	json.NewDecoder(r.Body).Decode(&newOrder)

	daysToBook := tools.DaysBetween(newOrder.From, newOrder.To)

	unavailableDays := make(map[time.Time]struct{})
	for _, day := range daysToBook {
		unavailableDays[day] = struct{}{}
	}

	for _, dayToBook := range daysToBook {
		for i, availability := range Availability {
			if !availability.Date.Equal(dayToBook) || availability.Quota < 1 {
				continue
			}
			availability.Quota -= 1
			Availability[i] = availability
			delete(unavailableDays, dayToBook)
		}
	}

	if len(unavailableDays) != 0 {
		http.Error(w, "Hotel room is not available for selected dates", http.StatusInternalServerError)
		utils.Logger.ErrorWithContext(r.Context())(ErrRoomNotAvailable, zap.Any("unavailableDays", unavailableDays), zap.Any("newOrder", newOrder))
		LogErrorf("Hotel room is not available for selected dates:\n%v\n%v", newOrder, unavailableDays)
		return
	}

	Orders = append(Orders, newOrder)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newOrder)

	LogInfo("Order successfully created: %v", newOrder)
}

var logger = log.Default()

func LogErrorf(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	logger.Printf("[Error]: %s\n", msg)
}

func LogInfo(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	logger.Printf("[Info]: %s\n", msg)
}
