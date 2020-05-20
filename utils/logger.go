package utils

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var onceInit sync.Once

//Logger global Logger
var Logger *zap.Logger

func init() {
	onceInit.Do(func() {
		initLogger()
	})
}

func initLogger() {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	config.OutputPaths = []string{"stdout"}
	Logger, _ = config.Build()
}
