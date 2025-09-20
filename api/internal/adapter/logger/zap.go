package logger

import (
	"context"
	"fmt"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"microservice/config"
	"microservice/pkg/utils"
	"net"
	"os"
	"time"
)

type level int

const (
	_ level = iota
	DEBUG
	INFO
	WARN
	ERR
	FATAL
)

type globalLogger struct {
	service   config.Service
	zap       zap.Logger
	logIngest net.Conn
}

func New(service config.Service, logIngest ILogIngest) IGLogger {
	gl := new(globalLogger)
	gl.service = service
	gl.logIngest = logIngest.Conn()
	return gl
}

func (gl *globalLogger) Init() {
	var (
		ns    = fmt.Sprintf("%s.%s", gl.service.NameSpace, gl.service.Name)
		cores = []zapcore.Core{stdoutCore(), fileCore()}
		opts  = make([]zap.Option, 0)
	)

	if gl.service.Debug == true {
		lscs := logstashCores(gl.logIngest)
		cores = append(cores, lscs[0], lscs[1])
	}

	opts = append(opts, zap.AddCaller(), zap.AddCallerSkip(1))

	// Combine the cores using zapcore.NewTee
	gl.zap = *zap.New(zapcore.NewTee(cores...), opts...)
	gl.zap = *gl.zap.With(
		zap.String("service.env", gl.service.Env),
		zap.String("service.namespace", ns),
	)
	gl.logIngest = nil
}

func (gl *globalLogger) C() *zap.Logger                          { return &gl.zap }
func (gl *globalLogger) Debug(scope string, fields ...zap.Field) { gl.zap.Debug(scope, fields...) }
func (gl *globalLogger) Info(scope string, fields ...zap.Field)  { gl.zap.Info(scope, fields...) }
func (gl *globalLogger) Warn(scope string, fields ...zap.Field)  { gl.zap.Warn(scope, fields...) }
func (gl *globalLogger) Error(scope string, fields ...zap.Field) { gl.zap.Error(scope, fields...) }

func (gl *globalLogger) Fx(lc fx.Lifecycle) ILogger {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) (err error) {
			utils.PrintStd(utils.StdLog, "logger", "[zap]initiated")
			return
		},
		OnStop: func(c context.Context) (err error) {
			utils.PrintStd(utils.StdLog, "logger", "[zap]stopping...")

			if gl.service.Debug == true {
				_ = gl.zap.Sync() // ignore Sync error: because the stdout isn't flushable
			} else {
				if err = gl.zap.Sync(); err != nil {
					utils.PrintStd(utils.StdLog, "logger", "[zap]failed to sync: %s", err)
					return
				}
			}

			utils.PrintStd(utils.StdLog, "logger", "[zap]stopped")
			return
		},
	})

	return gl
}

// HELPERS

func levelEnabler(lvl level) zap.LevelEnablerFunc {
	levels := map[level]zap.LevelEnablerFunc{
		DEBUG: zap.LevelEnablerFunc(debugLevel), // Logs everything to stdout
		INFO:  zap.LevelEnablerFunc(infoLevel),
		WARN:  zap.LevelEnablerFunc(warnLevel),
		ERR:   zap.LevelEnablerFunc(errorLevel),
	}

	return levels[lvl]
}

func debugLevel(lvl zapcore.Level) bool { return lvl >= zapcore.DebugLevel }

func infoLevel(lvl zapcore.Level) bool { return lvl == zapcore.InfoLevel }

func warnLevel(lvl zapcore.Level) bool { return lvl == zapcore.WarnLevel }

func errorLevel(lvl zapcore.Level) bool { return lvl == zapcore.ErrorLevel }

//

func stdoutCore() zapcore.Core {
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	stdoutWriter := zapcore.AddSync(os.Stdout)
	return zapcore.NewCore(consoleEncoder, stdoutWriter, levelEnabler(DEBUG))
}

//

func fileCore() zapcore.Core {
	fileName := fmt.Sprintf("%s/logs/%s.log", utils.Root(), time.Now().Format("2006-01-02"))
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	// NOTE: there could be multiple Cores per level

	return zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(file),
		levelEnabler(ERR),
	)
}

//

func logstashCores(conn net.Conn) [2]zapcore.Core {
	// Create an encoder for Logstash (JSON format)
	logstashEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	// Create a writer
	logstashWriter := zapcore.AddSync(conn)

	return [2]zapcore.Core{
		zapcore.NewCore(logstashEncoder, logstashWriter, levelEnabler(INFO)),
		zapcore.NewCore(logstashEncoder, logstashWriter, levelEnabler(WARN)),
	}
}
