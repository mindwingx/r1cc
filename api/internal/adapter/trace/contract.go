package trace

import (
	"context"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

type (
	IGTracer interface {
		Init()
		Fx(lc fx.Lifecycle) ITracer
		ITracer
	}

	ITracer interface {
		Provider() *sdktrace.TracerProvider
		SpanByCtx(c context.Context, operation, spType string, kind ...trace.SpanKind) (trace.Span, context.Context)
	}
)
