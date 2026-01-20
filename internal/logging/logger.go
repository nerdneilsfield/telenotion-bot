package logging

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Level string
	File  string
}

func NewLogger(cfg Config) (*zap.Logger, error) {
	level := zap.NewAtomicLevel()
	if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
	}

	config := zap.NewProductionConfig()
	config.Level = level
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}

	if cfg.File != "" {
		config.OutputPaths = append(config.OutputPaths, cfg.File)
		config.ErrorOutputPaths = append(config.ErrorOutputPaths, cfg.File)
	}

	config.EncoderConfig.TimeKey = "time"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return logger, nil
}
