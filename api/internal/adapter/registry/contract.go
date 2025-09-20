package registry

import "go.uber.org/fx"

type IRegistry interface {
	Init(configType, configFilePath string)
	ValueOf(key string) IRegistry
	Parse(item interface{}) error
	ParseNested(key string, item interface{}) error
	Fx(lc fx.Lifecycle) IRegistry
}
