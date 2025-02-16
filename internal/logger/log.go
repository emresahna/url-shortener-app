package logger

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func InitLogger() {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.TimeKey = "date"
	cfg.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder

	logger, err := cfg.Build()
	if err != nil {
		log.Fatalf("Error while initialize zap logger: %v", err)
	}

	Log = logger
}
