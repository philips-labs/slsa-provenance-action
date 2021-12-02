package options

import (
	"github.com/spf13/cobra"
)

// RootOptions Commandline flags used for the root command.
type RootOptions struct {
	Verbose bool
}

// AddFlags Registers the flags with the cobra.Command.
func (o *RootOptions) AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&o.Verbose, "verbose", "d", false, "show verbose output")
}
