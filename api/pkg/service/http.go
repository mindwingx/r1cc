package service

import (
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)

func Skipper(c echo.Context) (skip bool) {
	switch c.Path() {
	case "/handshake", "/metrics", "/public/swagger/*":
		skip = true
	}

	return
}

func OtelMiddleWare(serviceName string) echo.MiddlewareFunc {
	return otelecho.Middleware(serviceName, otelecho.WithSkipper(Skipper))
}
