package registry

import (
	"context"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"microservice/pkg/utils"
	"os"
	"path/filepath"
)

type registry struct {
	viper.Viper
}

func New() IRegistry {
	return &registry{*viper.New()}
}

func (r *registry) Init(configType, configFilePath string) {
	var err error

	switch configType {
	case ConfigTypeEnv:
		if err = r.parseEnvFile(configFilePath); err != nil {
			utils.PrintStd(utils.StdPanic, "registry", "init err: %s", err)
		}
	case ConfigTypeYAML, ConfigTypeYML:
		if err = r.parseYmlFile(configType, configFilePath); err != nil {
			utils.PrintStd(utils.StdPanic, "registry", "init err: %s", err)
		}
	default:
		utils.PrintStd(utils.StdPanic, "registry", "no config file found")
	}

	//NOTE: Todo - configs from remote services(like hashicorp-consul) will be placed here

	return
}

func (r *registry) ValueOf(key string) IRegistry {
	return &registry{*r.Sub(key)}
}

func (r *registry) Parse(item interface{}) (err error) {
	err = r.Unmarshal(&item)
	if err != nil {
		err = errors.New(fmt.Sprintf("failed to parse config file: %s", err.Error()))
		return err
	}

	return nil
}

func (r *registry) ParseNested(key string, item interface{}) (err error) {
	err = r.UnmarshalKey(key, &item)
	if err != nil {
		err = errors.New(fmt.Sprintf("failed to parse config file: %s", err.Error()))
		return err
	}

	return nil
}

func (r *registry) Fx(lc fx.Lifecycle) IRegistry {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) (err error) {
			utils.PrintStd(utils.StdLog, "registry", "initiated")
			return
		},
		OnStop: func(ctx context.Context) (err error) {
			utils.PrintStd(utils.StdLog, "registry", "stopped")
			return
		},
	})

	return r
}

// HELPERS

// parseEnvFile configFilePath accepts filename and it's path
func (r *registry) parseEnvFile(configFilePath string) (err error) {
	dir, file := filepath.Split(configFilePath)

	r.AddConfigPath(dir)
	r.SetConfigType(ConfigTypeEnv)
	r.SetConfigName(file)
	r.AutomaticEnv()

	if err = r.ReadInConfig(); err != nil {
		err = errors.New(fmt.Sprintf("failed to load config file: %s", err.Error()))
	}

	return
}

// parseYmlFile YAML and YML parser. the configFilePath accepts filename and it's path
func (r *registry) parseYmlFile(configType, configFilePath string) (err error) {
	r.SetConfigType(configType)

	config, err := os.Open(configFilePath)
	if err != nil {
		err = errors.New(fmt.Sprintf("failed to load config file: %s", err.Error()))
		return
	}

	if err = r.ReadConfig(config); err != nil {
		err = errors.New(fmt.Sprintf("failed to load config file: %s", err.Error()))
	}

	return
}
