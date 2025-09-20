package logger

import (
	"context"
	"fmt"
	"go.uber.org/fx"
	"microservice/config"
	"microservice/internal/adapter/registry"
	"microservice/pkg/utils"
	"net"
	"time"
)

type logger struct {
	service  config.Service
	config   config.Logger
	logstash net.Conn
}

func NewIngest(service config.Service, registry registry.IRegistry) ILog {
	l := new(logger)
	l.service = service

	if err := registry.Parse(&l.config); err != nil {
		utils.PrintStd(utils.StdPanic, "logger", "[logstash] parse err: %s", err)
	}

	return l
}

func (l *logger) Init() {
	conn, err := net.DialTimeout("tcp",
		fmt.Sprintf("%s:%s", l.config.Host, l.config.Port),
		time.Duration(l.config.Timeout)*time.Second,
	)

	if err != nil {
		utils.PrintStd(utils.StdPanic, "logger", "[logstash] init err: %s", err)
	}

	l.logstash = conn
}

func (l *logger) Conn() net.Conn {
	return l.logstash
}

func (l *logger) Write(p []byte) (res int, err error) {
	if l.service.Env == string(config.Dev) {
		res = 0
		return
	}

	return l.logstash.Write(p)
}

func (l *logger) Fx(lc fx.Lifecycle) ILogIngest {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) (err error) {
			utils.PrintStd(utils.StdLog, "logger", "[logstash]initiated")
			return
		},
		OnStop: func(ctx context.Context) (err error) {
			utils.PrintStd(utils.StdLog, "logger", "[logstash]stopping...")

			if err = l.logstash.Close(); err != nil {
				utils.PrintStd(utils.StdLog, "logger", "[logstash] connection close err: %s", err)
			}

			utils.PrintStd(utils.StdLog, "logger", "[logstash]stopped")
			return
		},
	})

	return l
}
