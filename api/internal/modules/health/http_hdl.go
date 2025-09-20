package health

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"microservice/internal/adapter/locale"
	"microservice/internal/adapter/trace"
	"microservice/pkg/meta"
	"microservice/pkg/meta/status"
	"time"
)

type (
	IHealthHttpHandler interface {
		Handshake(c echo.Context) error
	}

	HandlerFx struct {
		fx.In
		Locale locale.ILocale
		Tracer trace.ITracer
	}

	Handler struct {
		l   locale.ILocale
		trc trace.ITracer
	}

	Response struct {
		Message   string `json:"message"`
		Timestamp string `json:"timestamp"`
	}
)

func NewHttpHandlerFx(fx HandlerFx) IHealthHttpHandler {
	return &Handler{l: fx.Locale, trc: fx.Tracer}
}

// Handshake godoc
// @Summary Service Health
// @Description Aims to check the HTTP Service availability
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} Response
// @Router /handshake [get]
func (h *Handler) Handshake(c echo.Context) error {
	resp := Response{
		Message:   "connection established",
		Timestamp: time.Now().Format(time.DateTime),
	}

	return meta.Resp(c, h.l).Status(status.Success).Data(resp).Json()
}
