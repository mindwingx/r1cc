package middleware

import (
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"
	otelmtr "go.opentelemetry.io/otel/metric"
	"microservice/pkg/utils"
)

func (m *Middleware) RequestCounter(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		if c.Path() == "/metrics" || c.Path() == "/handshake" {
			err = next(c)
			return
		}

		counter, err := m.metric.Meter().Int64Counter(
			"http_requests_counter_general",
			otelmtr.WithDescription("http requests counter"),
		)

		if err != nil {
			utils.PrintStd(utils.StdLog, "metric", "[http] requests counter record err: %s", c.Path())
			err = next(c)
			return
		}

		err = next(c)

		counter.Add(c.Request().Context(), 1,
			otelmtr.WithAttributes(
				attribute.String("method", c.Request().Method),
				attribute.String("path", c.Path()),
				attribute.Int("status", c.Response().Status),
			),
		)

		return
	}
}
