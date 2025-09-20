package metric

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.uber.org/fx"
	"microservice/config"
	"microservice/pkg/utils"
	"os"
	"time"
)

type metrics struct {
	service  config.Service
	provider sdkmetric.MeterProvider
	meter    metric.Meter
}

func New(service config.Service) IMetrics {
	m := new(metrics)
	m.service = service
	return m
}

func (m *metrics) Init() {
	ctx := context.Background()

	exporter, err := prometheus.New()
	if err != nil {
		utils.PrintStd(utils.StdPanic, "metric", "failed to create OTLP exporter: %s", err)
	}

	name := fmt.Sprintf("%s.%s", m.service.NameSpace, m.service.Name)
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(name),
			semconv.ServiceNamespaceKey.String(m.service.NameSpace),
			semconv.ServiceInstanceIDKey.String(os.Getenv("HOSTNAME")), // HOSTNAME of kubernetes pod
			semconv.ServiceVersionKey.String(m.service.Version),
			semconv.DeploymentEnvironmentKey.String(m.service.Env),
		),
	)

	if err != nil {
		utils.PrintStd(utils.StdPanic, "metric", "failed to create resource: %s", err)
	}

	m.provider = *sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(exporter),
		sdkmetric.WithResource(res),
	)

	otel.SetMeterProvider(&m.provider) // Set the global meter provider
	m.meter = m.provider.Meter(name)
}

func (m *metrics) Fx(lc fx.Lifecycle) IMetric {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) (err error) {
			utils.PrintStd(utils.StdLog, "metric", "initiated")
			return
		},
		OnStop: func(ctx context.Context) (err error) {
			utils.PrintStd(utils.StdLog, "metric", "stopping...")

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err = m.provider.ForceFlush(ctx); err != nil {
				utils.PrintStd(utils.StdLog, "metric", "flush failed: %s", err.Error())
			}

			if err = m.provider.Shutdown(ctx); err != nil {
				utils.PrintStd(utils.StdLog, "metric", "shutdown error: %s", err.Error())
				return
			}

			utils.PrintStd(utils.StdLog, "metric", "stopped")
			return
		},
	})

	return m
}

// IMetric

func (m *metrics) Meter() metric.Meter {
	return m.meter
}
