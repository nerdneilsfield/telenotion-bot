package cmd

import (
	"fmt"

	"github.com/nerdneilsfield/telenotion-bot/internal/logging"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	verbose bool
	logger  *zap.Logger
)

func newRootCmd(version string, buildTime string, gitCommit string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "telenotion-bot",
		Short: "telenotion-bot captures Telegram messages into Notion.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			level := "info"
			if verbose {
				level = "debug"
			}
			newLogger, err := logging.NewLogger(logging.Config{Level: level})
			if err == nil {
				logger = newLogger
			}
		},
	}

	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")

	cmd.AddCommand(newVersionCmd(version, buildTime, gitCommit))
	cmd.AddCommand(newBotCmd())
	return cmd
}

func Execute(version string, buildTime string, gitCommit string) error {
	if logger == nil {
		baseLogger, err := logging.NewLogger(logging.Config{Level: "info"})
		if err == nil {
			logger = baseLogger
		}
	}

	if err := newRootCmd(version, buildTime, gitCommit).Execute(); err != nil {
		if logger != nil {
			logger.Error("error executing root command", zap.Error(err))
		}
		return fmt.Errorf("error executing root command: %w", err)
	}

	return nil
}
