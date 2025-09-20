package http

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"microservice/pkg/meta"
	"microservice/pkg/meta/status"
	"microservice/pkg/service"
)

func (s *Server) setMiddlewares() {
	// generic middleware
	s.client.Use(middleware.RequestID())
	s.client.Use(middleware.Secure())
	s.client.Use(middleware.BodyLimit(s.config.BodyLimit))

	if s.service.Debug == true {
		s.client.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format:  "method=${method}, uri=${uri}, status=${status}, ip=${ip} \n",
			Skipper: service.Skipper,
		}))
	}

	s.client.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         1 << 10, // 1 KB
		DisablePrintStack: false,
		DisableStackAll:   false,
		LogErrorFunc:      s.panicLogger,
	}))

	// custom middleware

	s.client.Use(
		s.middleware.RequestDuration,
		s.middleware.RequestCounter,
		s.middleware.RequestProcess,
	)
}

// HELPERS

func (s *Server) panicLogger(c echo.Context, err error, stack []byte) error {
	span, _ := s.trc.SpanByCtx(c.Request().Context(), "panic", "recover")
	defer span.End()

	span.RecordError(err, trace.WithAttributes(attribute.String("stack", string(stack))))
	s.lgr.Error("middleware.panic.logger", zap.Error(err))

	return meta.Resp(c, s.l).Status(status.Failed).Json()
}
