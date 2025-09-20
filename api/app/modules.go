package app

import (
	"go.uber.org/fx"
	"microservice/internal/modules/credit"
	"microservice/internal/modules/health"
	"microservice/internal/modules/message"
	"microservice/internal/modules/outbox"
	"microservice/internal/modules/tenant"
	"microservice/internal/modules/transaction"
)

type Modules []fx.Option

func (a *App) InitModules() {
	a.SetModule(&Modules{
		fx.Module("health", fx.Provide(health.NewHttpHandlerFx)),
		fx.Module("tenant", fx.Provide(tenant.NewRepositoryFx, tenant.NewUsecaseFx, tenant.NewHttpHandlerFx)),
		fx.Module("credit", fx.Provide(credit.NewRepositoryFx, credit.NewUsecaseFx, credit.NewHttpHandlerFx)),
		fx.Module("transaction", fx.Provide(transaction.NewRepositoryFx)),
		fx.Module("message", fx.Provide(message.NewRepositoryFx, message.NewUsecaseFx, message.NewHttpHandlerFx)),
		fx.Module("outbox", fx.Provide(outbox.NewRepositoryFx)),
	})

	a.Span().AddEvent("fx-modules initialized")
}
