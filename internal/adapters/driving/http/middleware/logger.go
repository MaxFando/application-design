package middlewares

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/MaxFando/application-design/pkg/utils"
)

func Logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		realIP := c.RealIP()

		reqBody, _ := io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(reqBody))

		respDumper := newResponseDumper(c.Response())
		responseStr := respDumper.GetResponse()
		timeStart := time.Now()

		utils.Logger.
			InfoWithContext(c.Request().Context())(
			"API",
			zap.String("url", c.Request().RequestURI),
			zap.String("body", string(reqBody)),
		)

		err := next(c)
		if err != nil {
			utils.Logger.
				InfoWithContext(c.Request().Context())(
				"API error",
				zap.String("IP", realIP),
				zap.String("method", req.Method),
				zap.String("path", c.Path()),
				zap.Any("request", string(reqBody)),
				zap.Error(err),
			)
		}

		timeEnd := time.Now()
		diffTime := timeEnd.Sub(timeStart)
		utils.Logger.
			InfoWithContext(c.Request().Context())(
			"API query",
			zap.String("IP", realIP),
			zap.String("method", req.Method),
			zap.String("path", c.Path()),
			zap.Any("request", string(reqBody)),
			zap.Any("time", diffTime),
			zap.Any("response", responseStr),
		)

		return err
	}
}

type responseDumper struct {
	http.ResponseWriter

	mw  io.Writer
	buf *bytes.Buffer
}

func newResponseDumper(resp *echo.Response) *responseDumper {
	buf := new(bytes.Buffer)
	return &responseDumper{
		ResponseWriter: resp.Writer,

		mw:  io.MultiWriter(resp.Writer, buf),
		buf: buf,
	}
}

func (d *responseDumper) GetResponse() string {
	return d.buf.String()
}
