package config

type Database struct {
	Debug              bool   `mapstructure:"-"`
	Host               string `mapstructure:"DB_HOST"`
	Port               string `mapstructure:"DB_PORT"`
	Username           string `mapstructure:"DB_USERNAME"`
	Password           string `mapstructure:"DB_PASSWORD"`
	Database           string `mapstructure:"DB_DATABASE"`
	Ssl                string `mapstructure:"DB_SSL"`
	MaxIdleConnections int    `mapstructure:"DB_MAX_IDLE_CONNECTIONS"`
	MaxOpenConnections int    `mapstructure:"DB_MAX_OPEN_CONNECTIONS"`
	MaxLifetimeSeconds int    `mapstructure:"DB_MAX_LIFETIME_SECONDS"`
	SlowSqlThreshold   int    `mapstructure:"DB_SLOW_SQL_THRESHOLD"`
}
