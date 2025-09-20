package metric

import (
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/fx"
)

type (
	IMetrics interface {
		Init()
		Fx(lc fx.Lifecycle) IMetric
		IMetric
	}

	IMetric interface {
		Meter() metric.Meter
	}
)
