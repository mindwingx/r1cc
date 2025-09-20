package trace

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"microservice/config"
	"microservice/internal/adapter/registry"
	"microservice/pkg/utils"
	"os"
	"time"
)

type tracer struct {
	service  config.Service
	config   config.Tracer
	exporter *otlptrace.Exporter
	provider *sdktrace.TracerProvider
}

func New(service config.Service, registry registry.IRegistry) IGTracer {
	t := new(tracer)
	t.service = service

	if err := registry.Parse(&t.config); err != nil {
		utils.PrintStd(utils.StdPanic, "tracer", "config parse err: %s", err)
	}

	return t
}

// IGTracer

func (t *tracer) Init() {
	ctx := context.Background()
	addr := fmt.Sprintf("%s:%s", t.config.Host, t.config.Port)

	exporter, err := httpExporter(ctx, addr)
	if err != nil {
		utils.PrintStd(utils.StdPanic, "tracer", "failed to create OTLP exporter: %s", err)
	}

	t.exporter = exporter

	name := fmt.Sprintf("%s.%s", t.service.NameSpace, t.service.Name)
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(name),
			semconv.ServiceNamespaceKey.String(t.service.NameSpace),
			semconv.ServiceInstanceIDKey.String(os.Getenv("HOSTNAME")), // HOSTNAME of kubernetes pod
			semconv.ServiceVersionKey.String(t.service.Version),
			semconv.DeploymentEnvironmentKey.String(t.service.Env),
		),
	)

	if err != nil {
		utils.PrintStd(utils.StdPanic, "tracer", "failed to create resource: %s", err)
	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	t.provider = provider
}

func (t *tracer) Provider() *sdktrace.TracerProvider {
	return t.provider
}

func (t *tracer) Fx(lc fx.Lifecycle) ITracer {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) (err error) {
			utils.PrintStd(utils.StdLog, "tracer", "initiated")
			return
		},
		OnStop: func(ctx context.Context) (err error) {
			utils.PrintStd(utils.StdLog, "tracer", "stopping...")

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err = t.provider.ForceFlush(ctx); err != nil {
				utils.PrintStd(utils.StdLog, "tracer", "flush failed: %s", err.Error())
			}

			if err = t.provider.Shutdown(ctx); err != nil {
				utils.PrintStd(utils.StdLog, "tracer", "shutdown error: %s", err.Error())
				return
			}

			utils.PrintStd(utils.StdLog, "tracer", "stopped")
			return
		},
	})

	return t
}

// ITracer

func (t *tracer) SpanByCtx(c context.Context, operation, spType string, kind ...trace.SpanKind) (trace.Span, context.Context) {
	spanName := fmt.Sprintf("%s.%s", operation, spType)
	spanKind := trace.SpanKindInternal

	if len(kind) > 0 && kind[0] != trace.SpanKindInternal {
		spanKind = kind[0]
	}

	ctx, span := otel.Tracer(t.service.Name).Start(c, spanName,
		trace.WithSpanKind(spanKind),
		trace.WithAttributes(attribute.String("span.type", spType)),
	)

	if !span.IsRecording() {
		utils.PrintStd(utils.StdLog, "tracer", "warning: inactive for request %s", spanName)
	}

	return span, ctx
}

// HELPERS

func httpExporter(ctx context.Context, otelEndpoint string) (*otlptrace.Exporter, error) {
	return otlptracehttp.New(
		ctx,
		otlptracehttp.WithEndpoint(otelEndpoint),
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithRetry(otlptracehttp.RetryConfig{
			Enabled:         true,
			InitialInterval: 1 * time.Second,
			MaxInterval:     5 * time.Second,
			MaxElapsedTime:  30 * time.Second,
		}),
	)
}
