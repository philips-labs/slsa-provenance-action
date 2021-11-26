package cli

import (
	"github.com/spf13/cobra"

	"github.com/philips-labs/slsa-provenance-action/cmd/slsa-provenance/cli/options"
)

const (
	cliName = "slsa-provenance"
)

var (
	ro = &options.RootOptions{}
)

// New creates a new instance of the slsa-provenance commandline interface
func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:               cliName,
		DisableAutoGenTag: true,
		SilenceUsage:      true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	ro.AddFlags(cmd)

	cmd.AddCommand(Version())
	cmd.AddCommand(Files())

	return cmd
}
