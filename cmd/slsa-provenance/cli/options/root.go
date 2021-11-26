package options

import (
	"github.com/spf13/cobra"
)

type RootOptions struct {
	Verbose bool
}

func (o *RootOptions) AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&o.Verbose, "verbose", "d", false, "show verbose output")
}
