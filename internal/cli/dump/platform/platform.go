// Package platform implements the platform dump subcommand for godyl.
// It displays information about the detected system platform.
package platform

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/flags"
	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/detect"
	iutils "github.com/idelchi/godyl/internal/utils"
)

// Command encapsulates the platform dump command with its associated configuration.
type Command struct {
	// Command is the platform cobra.Command instance
	Command *cobra.Command
}

// Flags adds platform-specific flags to the command.
func (cmd *Command) Flags() {
	// No specific flags for this command
}

// NewPlatformCommand creates a Command for displaying platform information.
func NewPlatformCommand(cfg *config.Config) *Command {
	cmd := &cobra.Command{
		Use:   "platform",
		Short: "Display platform information",
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			return flags.ChainPreRun(cmd, nil)
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			c, err := getPlatform()
			if err != nil {
				return err
			}

			iutils.Print(cfg.Dump.Format, c)

			return nil
		},
	}

	return &Command{
		Command: cmd,
	}
}

// NewCommand creates a cobra.Command instance for the platform dump subcommand.
func NewCommand(cfg *config.Config) *cobra.Command {
	// Create the platform command
	cmd := NewPlatformCommand(cfg)

	// Add platform-specific flags
	cmd.Flags()

	return cmd.Command
}

// getPlatform detects and returns information about the current platform.
func getPlatform() (*detect.Platform, error) {
	platform := &detect.Platform{}
	if err := platform.Detect(); err != nil {
		return nil, fmt.Errorf("detecting platform: %w", err)
	}

	return platform, nil
}
