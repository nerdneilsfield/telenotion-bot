package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/nerdneilsfield/telenotion-bot/cmd"
	"github.com/nerdneilsfield/telenotion-bot/internal/logging"
	"go.uber.org/zap"
)

var (
	version   = "dev"
	buildTime = "unknown"
	gitCommit = "unknown"
)

var logger *zap.Logger

func init() {
	baseLogger, err := logging.NewLogger(logging.Config{Level: "info"})
	if err == nil {
		logger = baseLogger
	}
}

func gracefulShutdown() {
	if logger != nil {
		logger.Info("Shutting down...")
		logger.Sync()
	}
}

func main() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalChan
		gracefulShutdown()
		os.Exit(0)
	}()

	if err := cmd.Execute(version, buildTime, gitCommit); err != nil {
		if logger != nil {
			logger.Error("Failed to execute root command", zap.Error(err))
		}
		os.Exit(1)
	}
}
