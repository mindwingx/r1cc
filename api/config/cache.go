package config

import "time"

type Cache struct {
	DB       int           `mapstructure:"CACHE_DB"`
	Host     string        `mapstructure:"CACHE_HOST"`
	Port     string        `mapstructure:"CACHE_PORT"`
	Username string        `mapstructure:"CACHE_USERNAME"`
	Password string        `mapstructure:"CACHE_PASSWORD"`
	Timeout  time.Duration `mapstructure:"CACHE_TIMEOUT"`
}
