package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
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

	cmd.PersistentFlags().StringVarP(&configPath, "config", "c", "config.toml", "Path to configuration file")

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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalChan
		cancel()
	}()

	runners := make([]func(context.Context) error, 0, 2)
	if cfg.Telegram.Token != "" {
		telegramRunner, err := session.NewRunner(cfg, logger)
		if err != nil {
			return err
		}
		runners = append(runners, telegramRunner.Run)
	}
	if cfg.Discord.Token != "" {
		discordRunner, err := session.NewDiscordRunner(cfg, logger)
		if err != nil {
			return err
		}
		runners = append(runners, discordRunner.Run)
	}

	var wg sync.WaitGroup
	errCh := make(chan error, len(runners))

	for _, run := range runners {
		wg.Add(1)
		go func(run func(context.Context) error) {
			defer wg.Done()
			if err := run(ctx); err != nil {
				errCh <- err
			}
		}(run)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	select {
	case err := <-errCh:
		if err != nil {
			cancel()
			return err
		}
	case <-ctx.Done():
		return nil
	}

	return nil
}
