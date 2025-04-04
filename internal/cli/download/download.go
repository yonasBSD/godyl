// Package download implements the download command for godyl.
// It provides functionality to download and extract tools from various sources.
package download

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/flags"
	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/core/defaults"
	"github.com/idelchi/godyl/internal/core/processor"
	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/internal/tools/sources"
	iutils "github.com/idelchi/godyl/internal/utils"
	"github.com/idelchi/godyl/pkg/logger"
	"github.com/idelchi/godyl/pkg/utils"
)

// Command encapsulates the download cobra command with its associated config and embedded files.
type Command struct {
	// Command is the download cobra.Command instance
	Command *cobra.Command
}

// Flags adds download-specific flags to the command.
func (cmd *Command) Flags() {
	flags.Tool(cmd.Command)
}

// NewDownloadCommand creates a Command for downloading and unpacking tools.
func NewDownloadCommand(cfg *config.Config, files config.Embedded) *Command {
	// Create the download command
	cmd := &cobra.Command{
		Use:     "download [tool]",
		Aliases: []string{"dl", "unpack", "extract", "x"},
		Short:   "Download and unpack tools",
		Long:    "Download and unpack tools from GitHub, URLs, or Go projects",
		Args:    cobra.MinimumNArgs(1),
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			return flags.ChainPreRun(cmd, &cfg.Tool, cmd.Root().Name(), "tool")
		},
		RunE: func(_ *cobra.Command, args []string) error {
			lvl, err := logger.LevelString(cfg.Root.Log)
			if err != nil {
				return fmt.Errorf("parsing log level: %w", err)
			}

			log := logger.New(lvl)

			// Load defaults
			defaults, err := defaults.Load(cfg.Root.Defaults, files, *cfg)
			if err != nil {
				return fmt.Errorf("loading defaults: %w", err)
			}

			if cfg.Tool.Show {
				iutils.Print("yaml", cfg.Root, cfg.Tool)

				return nil
			}

			toolsList := []tools.Tool{}

			var version string
			utils.SetIfZeroValue(&version, cfg.Tool.Version)

			for _, name := range args {
				tool := tools.Tool{
					Name: name,
					Mode: tools.Extract,
					Version: tools.Version{
						Version: version,
					},
				}
				if utils.IsURL(name) {
					tool.Name = filepath.Base(name)
					tool.Path = name
					tool.Source.Type = sources.DIRECT
					tool.Version.Version = version
				}

				toolsList = append(toolsList, tool)
			}

			// Process tools
			proc := processor.New(toolsList, defaults, *cfg, log)
			if err := proc.Process(nil, nil); err != nil {
				return fmt.Errorf("processing tools: %w", err)
			}

			return nil
		},
	}

	return &Command{
		Command: cmd,
	}
}

// NewCommand creates a cobra.Command instance containing the download command.
func NewCommand(cfg *config.Config, files config.Embedded) *cobra.Command {
	// Create the download command
	cmd := NewDownloadCommand(cfg, files)

	// Add tool-specific flags
	cmd.Flags()

	return cmd.Command
}
