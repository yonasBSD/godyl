package sync

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cache"
	"github.com/idelchi/godyl/internal/cli/flags"
	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/defaults"
	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/internal/tools/tags"
	"github.com/idelchi/godyl/internal/utils"
	"github.com/idelchi/godyl/pkg/executable"
	"github.com/idelchi/godyl/pkg/logger"
	"github.com/idelchi/godyl/pkg/path/files"
	"github.com/idelchi/godyl/pkg/version"
)

type Command struct {
	// Command is the Cache cobra.Command instance
	Command *cobra.Command
}

func (cmd *Command) Flags() {
}

func NewSyncCommand(cfg *config.Config, embedded config.Embedded) *Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Interact with the cache",
		Long:  "Interact with the cache.",
		Args:  cobra.ArbitraryArgs,

		PreRunE: func(cmd *cobra.Command, _ []string) error {
			return flags.ChainPreRun(cmd, nil)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			lvl, err := logger.LevelString(cfg.Root.Log)
			if err != nil {
				return fmt.Errorf("parsing log level: %w", err)
			}

			// Set the tools file if provided as an argument
			if len(args) > 0 {
				cfg.Tool.Tools = files.New("", args...)
			} else {
				cfg.Tool.Tools = files.New(".", "tools.yml")
			}

			log := logger.New(lvl)

			// Load defaults
			defaults, err := defaults.Load(cfg.Root.Defaults, embedded, *cfg)
			if err != nil {
				return fmt.Errorf("loading defaults: %w", err)
			}

			toolsList := tools.Tools{}

			// Load tools
			for _, file := range cfg.Tool.Tools {
				tools, err := utils.LoadTools(file, defaults, cfg.Root.Default)
				if err != nil {
					return fmt.Errorf("loading tools: %w", err)
				}

				toolsList = append(toolsList, tools...)
			}

			cacheHandler := cache.New(cfg.Root.Cache.Dir)
			if err := cacheHandler.Load(); err != nil {
				return fmt.Errorf("loading cache: %w", err)
			}

			for _, tool := range toolsList {
				if res := tool.Resolve(tags.IncludeTags{}, true); res.IsFailed() {
					return fmt.Errorf("resolving tool: %s", res)
				}

				// Parse the version of the existing tool.
				exe := executable.New(tool.AbsPath())

				id := tool.ID()
				commands := tool.Version

				item, err := cacheHandler.Get(id)

				if !exe.Exists() {
					if !errors.Is(err, cache.ErrItemNotFound) {
						if err := cacheHandler.Delete(id); err != nil {
							log.Warn("failed to delete cache: %v", err)
						} else {
							log.Warn("cache deleted for %q: executable %q has been removed from system", item.Name, item.Path)
						}
					}

					continue
				}

				if err != nil {
					log.Warn("failed to get cache for %q: %v", tool.Name, err)

					continue
				}

				// Check if we have commands to determine version
				if commands.Commands == nil || len(commands.Commands) == 0 {
					continue
				}

				// Parse version using available commands
				parser := &executable.Parser{
					Patterns: commands.Patterns,
					Commands: commands.Commands,
				}

				parsed, err := exe.Parse(parser)
				if err != nil {
					log.Warn("failed to parse version for %q: %v", item.Name, err)

					continue
				}

				if version.Compare(parsed, item.Version) {
					continue
				}

				item.Version = parsed
				item.Updated = time.Now()

				if err := cacheHandler.Save(item); err != nil {
					log.Warn("failed to save cache for %q: %v", item.Name, err)
				} else {
					log.Info("cache updated for %q: version %q parsed", item.Name, item.Version)
				}
			}

			if !cacheHandler.Touched() {
				log.Info("no changes necessary")
			}

			return nil
		},
	}

	return &Command{
		Command: cmd,
	}
}

func NewCommand(cfg *config.Config, files config.Embedded) *cobra.Command {
	cmd := NewSyncCommand(cfg, files)

	// Add Cache-specific flags
	cmd.Flags()

	return cmd.Command
}
