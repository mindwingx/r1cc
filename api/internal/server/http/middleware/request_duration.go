package middleware

import (
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"
	otelmtr "go.opentelemetry.io/otel/metric"
	"microservice/pkg/utils"
	"time"
)

func (m *Middleware) RequestDuration(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		if c.Path() == "/handshake" || c.Path() == "/metrics" {
			err = next(c)
			return
		}

		duration, err := m.metric.Meter().Float64Histogram(
			"http_requests_duration_general",
			otelmtr.WithDescription("http requests duration"),
			otelmtr.WithUnit("ms"),
		)

		if err != nil {
			utils.PrintStd(utils.StdLog, "metric", "[http] requests duration record err: %s", c.Path())
			err = next(c)
			return
		}

		start := time.Now()
		err = next(c)

		// record the duration after the "next", receives the status code

		d := time.Since(start).Seconds()

		duration.Record(c.Request().Context(), d,
			otelmtr.WithAttributes(
				attribute.String("method", c.Request().Method),
				attribute.String("path", c.Path()),
				attribute.Int("status", c.Response().Status),
				attribute.Float64("duration", d,
				),
			),
		)

		return
	}
}
