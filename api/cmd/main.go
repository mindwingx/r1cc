package main

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"microservice/app"
	"microservice/internal/server/http/middleware"
	"microservice/pkg/utils"
	"os"
	"time"
)

func main() {
	st, ln := utils.Preload()

	service := app.New()
	{
		service.InitClients()
		service.InitProviders()
		service.InitModules()
		service.InitServers()
	}

	var options []fx.Option
	{
		options = append(options, *service.Module()...)

		options = append(options,
			fx.Provide(service.Client().Registry().Fx),
			fx.Provide(service.Client().Tracer().Fx),
			fx.Provide(service.Client().Logger().Fx),
			fx.Provide(service.Client().GlobalLogger().Fx),
			fx.Provide(service.Client().Metric().Fx),
			fx.Provide(service.Client().Locale().Fx),
			fx.Provide(service.Client().Cache().Fx),
			fx.Provide(service.Client().Database().Fx),
			fx.Provide(service.Client().Queue().Fx),
		)

		options = append(options, *service.Provider()...)
		options = append(options,
			fx.Provide(middleware.NewFx),
			fx.Invoke(service.Server().Http().Fx),
		)
	}

	fxNew := fx.New(options...)
	{
		service.Span().AddEvent("service initialized",
			trace.WithAttributes(attribute.String("duration", time.Since(st).String())),
		)

		startSp, _ := service.Client().Tracer().SpanByCtx(service.Ctx(), "service", "start")

		ctxx := context.Background()
		ctxs := context.WithValue(ctxx, "mwx", "nex")

		if err := fxNew.Start(ctxs); err != nil {
			msg := fmt.Errorf("start error: %s", err)
			startSp.RecordError(msg)
			utils.PrintStd(utils.StdPanic, "service err: %s", msg.Error())
		}

		startSp.AddEvent("started")
		func() {
			defer service.Span().End()
			defer startSp.End()
		}()

		if val := <-fxNew.Done(); val == os.Interrupt || val == os.Kill {
			utils.PrintStd(utils.Std, "", "\n%s [service] shutting down %s\n", ln, ln)

			stopCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err := fxNew.Stop(stopCtx); err != nil {
				msg := fmt.Errorf("shutdown error: %s", err)
				utils.PrintStd(utils.StdLog, "service", msg.Error())
				return
			}

			return
		}
	}
}
