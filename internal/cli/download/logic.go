package download

import (
	"fmt"

	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/processor"
	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/internal/tools/hints"
	"github.com/idelchi/godyl/internal/tools/mode"
	"github.com/idelchi/godyl/internal/tools/sources"
	"github.com/idelchi/godyl/internal/tools/strategy"
	"github.com/idelchi/godyl/internal/tools/tags"
	"github.com/idelchi/godyl/internal/tools/tool"
	"github.com/idelchi/godyl/internal/tools/version"
	"github.com/idelchi/godyl/pkg/path/file"
	"github.com/idelchi/godyl/pkg/utils"
)

// run executes the `download` command.
func run(input common.Input) error {
	cfg, embedded, _, _, args := input.Unpack()

	if cfg.Download.Dry {
		cfg.Verbose = 1
	}

	tools := tools.Tools{}

	for _, name := range args {
		tool := &tool.Tool{
			Mode:     mode.Extract,
			Strategy: strategy.Force,
			Version: version.Version{
				Version: cfg.Download.Version,
			},
		}

		if utils.IsURL(name) {
			tool.Name = file.New(name).Base()
			tool.URL = name
			tool.Source.Type = sources.URL
		} else {
			tool.Name = name
			tool.Source.Type = sources.GITHUB
		}

		tools.Append(tool)
	}

	// Generate a common configuration for the command
	cfg.Common = cfg.Download.ToCommon()

	cfg.Cache.Disabled = true

	runner := common.NewHandler(*cfg, *embedded)
	if err := runner.SetupLogger(cfg.LogLevel); err != nil {
		return fmt.Errorf("setting up logger: %w", err)
	}

	if err := runner.Resolve(cfg.Defaults, &tools); err != nil {
		return err
	}

	// Add the hints that were passed via the `--hint` flag
	for _, tool := range tools {
		for _, hint := range cfg.Download.Hints {
			tool.Hints.Add(hints.Hint{
				Pattern: hint,
			})
		}
	}

	// Process tools
	proc := processor.New(tools, *cfg, runner.Logger())

	proc.NoDownload = cfg.Download.Dry
	if err := proc.Process(tags.IncludeTags{}); err != nil {
		return fmt.Errorf("processing tools: %w", err)
	}

	return nil
}
