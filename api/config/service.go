package config

type Env string

const (
	Prod  Env = "production"
	Stage Env = "stage"
	Dev   Env = "development"
)

type Service struct {
	//Locale    `mapstructure:",squash"` // squash is used to flatten the embedded structs for Viper
	Name      string `mapstructure:"APP_NAME"`
	Version   string `mapstructure:"APP_VERSION"`
	Env       string `mapstructure:"APP_ENV"`
	NameSpace string `mapstructure:"APP_NAMESPACE"`
	Debug     bool   `mapstructure:"APP_DEBUG"`
	TimeZone  string `mapstructure:"APP_TIMEZONE"`
}
