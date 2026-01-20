package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/nerdneilsfield/telenotion-bot/internal/config"
	"github.com/nerdneilsfield/telenotion-bot/internal/logging"
	"github.com/nerdneilsfield/telenotion-bot/internal/session"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var configPath string

func newBotCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bot",
		Short: "Run the Telegram to Notion capture bot",
		RunE:  runBot,
	}

	cmd.PersistentFlags().StringVar(&configPath, "config", "config.toml", "Path to configuration file")

	return cmd
}

func runBot(cmd *cobra.Command, args []string) error {
	if configPath == "" {
		return fmt.Errorf("--config is required")
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	configuredLogger, err := logging.NewLogger(logging.Config{Level: cfg.Log.Level, File: cfg.Log.File})
	if err != nil {
		return err
	}
	logger = configuredLogger
	defer logger.Sync()

	logger.Info("Starting bot", zap.String("config", configPath))

	runner, err := session.NewRunner(cfg, logger)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalChan
		cancel()
	}()

	return runner.Run(ctx)
}
