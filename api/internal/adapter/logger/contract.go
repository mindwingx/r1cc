package logger

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net"
)

type (
	ILog interface {
		Init()
		Fx(lc fx.Lifecycle) ILogIngest
		ILogIngest
	}

	ILogIngest interface {
		Conn() net.Conn
		Write(p []byte) (int, error)
	}
)

type (
	IGLogger interface {
		Init()
		Fx(lc fx.Lifecycle) ILogger
		ILogger
	}

	ILogger interface {
		C() *zap.Logger
		Debug(scope string, fields ...zap.Field)
		Info(scope string, fields ...zap.Field)
		Warn(scope string, fields ...zap.Field)
		Error(scope string, fields ...zap.Field)
	}
)
