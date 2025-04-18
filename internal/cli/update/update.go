// Package update implements the update command for godyl.
// It provides functionality to update the application itself to the latest version.
package update

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cache"
	"github.com/idelchi/godyl/internal/cli/flags"
	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/core/defaults"
	"github.com/idelchi/godyl/internal/core/updater"
	"github.com/idelchi/godyl/pkg/logger"
)

// Command encapsulates the update cobra command with its associated config and embedded files.
type Command struct {
	// Command is the update cobra.Command instance
	Command *cobra.Command
}

// Flags adds update-specific flags to the command.
func (cmd *Command) Flags() {
	flags.Update(cmd.Command)
}

// NewUpdateCommand creates a Command for updating the application to the latest version.
func NewUpdateCommand(cfg *config.Config, embedded config.Embedded) *Command {
	cmd := &cobra.Command{
		Use:     "update",
		Aliases: []string{"upgrade", "u"},
		Short:   "Update the application",
		Long:    "Update the godyl application to the latest version",
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			return flags.ChainPreRun(cmd, &cfg.Update)
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			lvl, err := logger.LevelString(cfg.Root.Log)
			if err != nil {
				return fmt.Errorf("parsing log level: %w", err)
			}

			defaults, err := defaults.Load(cfg.Root.Defaults, embedded, *cfg)
			if err != nil {
				return fmt.Errorf("loading defaults: %w", err)
			}
			cache, err := cache.New(cfg.Root.Cache.Dir)
			if err != nil {
				return fmt.Errorf("creating cache: %w", err)
			}

			appUpdater := updater.New(defaults, cfg.Update.NoVerifySSL, embedded.Template, lvl, cache)

			versions := updater.Versions{
				Current:   cmd.Root().Version,
				Requested: cfg.Update.Version,
				Pre:       cfg.Update.Pre,
			}

			if err := appUpdater.Update(versions); err != nil {
				return fmt.Errorf("updating: %w", err)
			}

			return nil
		},
	}

	return &Command{
		Command: cmd,
	}
}

// NewCommand creates a cobra.Command instance containing the update command.
func NewCommand(cfg *config.Config, files config.Embedded) *cobra.Command {
	// Create the update command
	cmd := NewUpdateCommand(cfg, files)

	// Add update-specific flags
	cmd.Flags()

	return cmd.Command
}
