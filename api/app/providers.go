package app

import (
	"go.uber.org/fx"
	"microservice/internal/adapter/provider/sms"
)

type Providers []fx.Option

func (a *App) InitProviders() {
	a.SetProvider(&Providers{
		fx.Module("provider.sms", fx.Provide(sms.New)),
	})
}
