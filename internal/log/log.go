package log

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var once sync.Once

var logger *zap.Logger

func Get() *zap.Logger {

	once.Do(func() {
		encoderCfg := zap.NewProductionEncoderConfig()
		encoderCfg.TimeKey = "timestamp"
		encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
		c := zap.Config{
			Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
			Development:       false,
			DisableCaller:     false,
			DisableStacktrace: false,
			Sampling:          nil,
			Encoding:          "json",
			EncoderConfig:     encoderCfg,
			OutputPaths: []string{
				"stderr",
			},
			ErrorOutputPaths: []string{
				"stderr",
			},
			InitialFields: map[string]interface{}{
				"pid": os.Getpid(),
			},
		}
		logger = zap.Must(c.Build())
	})

	return logger

}
