package middleware

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"microservice/config"
	"microservice/internal/adapter/cache"
	"microservice/internal/adapter/locale"
	"microservice/internal/adapter/logger"
	"microservice/internal/adapter/metric"
)

type MiddlewaresFx struct {
	fx.In
	Locale locale.ILocale
	Logger logger.ILogger
	Cache  cache.ICache
	Metric metric.IMetric
}
type Middleware struct {
	l      locale.ILocale
	lgr    logger.ILogger
	cache  cache.ICache
	metric metric.IMetric
	//
	service *config.Service
	router  *echo.Router
}

func NewFx(fx MiddlewaresFx) IMiddleware {
	return &Middleware{
		l:      fx.Locale,
		lgr:    fx.Logger,
		cache:  fx.Cache,
		metric: fx.Metric,
	}
}

func (m *Middleware) Service() *config.Service {
	return m.service
}
