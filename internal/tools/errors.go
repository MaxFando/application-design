package tools

import (
	"github.com/labstack/echo/v4"
)

type errorsMessage struct {
	StatusCode int    `json:"status"`
	Message    string `json:"message"`
}

func SendEchoHttpError(code int, message string) *echo.HTTPError {
	return echo.NewHTTPError(
		code,
		errorsMessage{
			StatusCode: code,
			Message:    message,
		},
	)
}
