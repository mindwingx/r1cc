package config

import "time"

type HTTP struct {
	Host         string        `mapstructure:"HTTP_SERVER_HOST" required:"true"`
	Port         string        `mapstructure:"HTTP_SERVER_PORT" required:"true"`
	WriteTimeout time.Duration `mapstructure:"HTTP_SERVER_WRITE_TIMEOUT"`
	ReadTimeout  time.Duration `mapstructure:"HTTP_SERVER_READ_TIMEOUT"`
	Tls          bool          `mapstructure:"HTTP_SERVER_TLS"`
	BodyLimit    string        `mapstructure:"HTTP_SERVER_BODY_LIMIT"`
}
