package config

type Logger struct {
	Host    string `mapstructure:"LOGGER_HOST" validate:"required"`
	Port    string `mapstructure:"LOGGER_PORT" validate:"required"`
	Timeout int    `mapstructure:"LOGGER_TIMEOUT"`
}
