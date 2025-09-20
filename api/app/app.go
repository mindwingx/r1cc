package app

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"microservice/config"
	"microservice/pkg/utils"
	_ "microservice/pkg/validator"
)

type (
	App struct {
		config   *config.Service
		client   *Clients
		provider *Providers
		module   *Modules
		server   *Servers
		trc      Tracer
	}

	Tracer struct {
		sp  trace.Span
		ctx context.Context
	}
)

func New() *App { return &App{} }

func (a *App) Config() *config.Service          { return a.config }
func (a *App) SetConfig(config *config.Service) { a.config = config }

func (a *App) Client() *Clients          { return a.client }
func (a *App) SetClient(client *Clients) { a.client = client }

func (a *App) Provider() *Providers            { return a.provider }
func (a *App) SetProvider(provider *Providers) { a.provider = provider }

func (a *App) Module() *Modules          { return a.module }
func (a *App) SetModule(module *Modules) { a.module = module }

func (a *App) Server() *Servers          { return a.server }
func (a *App) SetServer(server *Servers) { a.server = server }

func (a *App) Trace() *Tracer                                { return &a.trc }
func (a *App) Span() trace.Span                              { return a.trc.sp }
func (a *App) Ctx() context.Context                          { return a.trc.ctx }
func (a *App) SetTrace(span trace.Span, ctx context.Context) { a.trc.sp = span; a.trc.ctx = ctx }

// init

func (a *App) initService() {
	if err := a.Client().Registry().Parse(&a.config); err != nil {
		utils.PrintStd(utils.StdPanic, "service", "config parse err: %s", err)
	}
}
