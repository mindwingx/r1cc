package middleware

import (
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"
	otelmtr "go.opentelemetry.io/otel/metric"
	"microservice/pkg/utils"
)

func (m *Middleware) RequestProcess(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		if c.Path() == "/metrics" || c.Path() == "/handshake" {
			err = next(c)
			return
		}

		processCounter, err := m.metric.Meter().Int64Gauge(
			"http_requests_in_progress_counter_general",
			otelmtr.WithDescription("http requests process(in progress) counter"),
		)

		if err != nil {
			utils.PrintStd(utils.StdLog, "metric", "[http] requests process counter record err: %s", c.Path())
			err = next(c)
			return
		}

		attrs := []attribute.KeyValue{
			attribute.String("method", c.Request().Method),
			attribute.String("path", c.Path()),
		}

		processCounter.Record(c.Request().Context(), 1, otelmtr.WithAttributes(attrs...))

		err = next(c)
		processCounter.Record(c.Request().Context(), -1, otelmtr.WithAttributes(attrs...))

		return
	}
}
