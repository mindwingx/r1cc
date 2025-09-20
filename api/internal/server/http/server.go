package http

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"microservice/config"
	"microservice/internal/adapter/cache"
	"microservice/internal/adapter/locale"
	"microservice/internal/adapter/logger"
	"microservice/internal/adapter/metric"
	"microservice/internal/adapter/registry"
	"microservice/internal/adapter/trace"
	"microservice/internal/modules/credit"
	"microservice/internal/modules/health"
	"microservice/internal/modules/message"
	"microservice/internal/modules/tenant"
	"microservice/internal/server/http/middleware"
	"microservice/pkg/utils"
	"net/http"
	"time"
)

type (
	ServerFx struct {
		fx.In
		Locale     locale.ILocale
		Logger     logger.ILogger
		Tracer     trace.ITracer
		Metric     metric.IMetric
		Cache      cache.ICache
		Middleware middleware.IMiddleware
		//
		Health  health.IHealthHttpHandler
		Tenant  tenant.ITenantHttpHandler
		Credit  credit.ICreditHttpHandler
		Message message.IMessageHttpHandler
	}

	Server struct {
		*Handler
		trc        trace.ITracer
		mtr        metric.IMetric
		l          locale.ILocale
		lgr        logger.ILogger
		cache      cache.ICache
		middleware middleware.IMiddleware
		service    *config.Service
		config     *config.HTTP
		swagger    *config.Swagger
		client     *echo.Echo
	}

	Handler struct {
		health  health.IHealthHttpHandler
		tenant  tenant.ITenantHttpHandler
		credit  credit.ICreditHttpHandler
		message message.IMessageHttpHandler
	}
)

func NewServer(service *config.Service, registry registry.IRegistry) IHttpServer {
	s := new(Server)

	if err := registry.Parse(&s.config); err != nil {
		utils.PrintStd(utils.StdPanic, "http", "config parse err: %s", err)
	}

	if err := registry.Parse(&s.swagger); err != nil {
		utils.PrintStd(utils.StdPanic, "http", "swagger config parse err: %s", err)
	}

	host := s.config.Host
	if service.Env == string(config.Dev) {
		host = "localhost"
	}

	s.swagger.Host = fmt.Sprintf("%s:%s", host, s.config.Port)
	s.service = service

	return s
}

func (s *Server) Init() {
	s.client = echo.New()
}

func (s *Server) Fx(lc fx.Lifecycle, sfx ServerFx) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) (err error) {
			utils.PrintStd(utils.StdLog, "http", "initiated")
			defer utils.PrintStd(utils.Std, "http", "server is listening on port %s", s.config.Port)

			//note: do not change the order of the items below. they are not independent.
			s.l = sfx.Locale
			s.lgr = sfx.Logger
			s.trc = sfx.Tracer
			s.mtr = sfx.Metric
			s.cache = sfx.Cache
			s.middleware = sfx.Middleware
			s.Handler = &Handler{
				health:  sfx.Health,
				tenant:  sfx.Tenant,
				credit:  sfx.Credit,
				message: sfx.Message,
			}

			s.setupServer()

			if err = s.run(); err == nil {
				if s.service.Debug == true {
					s.routesStdout()
				}
			}

			return
		},
		OnStop: func(ctx context.Context) (err error) {
			utils.PrintStd(utils.StdLog, "http", "stopping...")
			defer utils.PrintStd(utils.StdLog, "http", "stopped")

			err = s.shutdown()
			return
		},
	})
}

//

func (s *Server) setupServer() {
	s.client.Debug = s.service.Debug
	s.setMiddlewares()
	s.setRoutes()
}

func (s *Server) run() (err error) {
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", s.config.Host, s.config.Port),
		ReadTimeout:  s.config.ReadTimeout,
		WriteTimeout: s.config.WriteTimeout,
	}

	go func() {
		if err = s.client.StartServer(server); err != nil {
			utils.PrintStd(utils.StdPanic, "http", "server start failure: %s", err)
		}
	}()

	return
}

func (s *Server) shutdown() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err = s.client.Server.Shutdown(ctx); err != nil {
		utils.PrintStd(utils.StdLog, "http", "server shutting  down", err)
	}

	return
}
