// Package download contains the subcommand definition for `download`.
package download

import (
	"github.com/spf13/cobra"

	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/config/download"
	"github.com/idelchi/godyl/internal/config/root"
)

// Command returns the `download` command.
func Command(global *root.Config, local any, embedded *common.Embedded) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "download [tool]",
		Aliases: []string{"dl", "x"},
		Short:   "Download and extract tools",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Exit early if the command is run with `--show/-s` flag.
			if common.ExitOnShow(global.ShowFunc) {
				return nil
			}

			return run(common.Input{Global: global, Cmd: cmd, Args: args, Embedded: embedded})
		},
	}

	common.SetSubcommandDefaults(cmd, local, global.ShowFunc)

	download.Flags(cmd)

	return cmd
}
