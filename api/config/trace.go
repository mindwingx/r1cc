package config

type Tracer struct {
	Host     string `mapstructure:"TRACE_HOST" validate:"required"`
	Port     string `mapstructure:"TRACE_PORT" validate:"required"`
	LogSpans bool   `mapstructure:"TRACE_LOG_SPANS"`
}
