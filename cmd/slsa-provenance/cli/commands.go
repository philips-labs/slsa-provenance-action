package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/philips-labs/slsa-provenance-action/cmd/slsa-provenance/cli/options"
)

const (
	cliName = "slsa-provenance"
)

var (
	ro = &options.RootOptions{}
)

// RequiredFlagError creates a required flag error for the given flag name
func RequiredFlagError(flagName string) error {
	return fmt.Errorf("no value found for required flag: %s", flagName)
}

// RequiredEnvironmentVariableError creates a required environment variable
// error for the given environment variable name
func RequiredEnvironmentVariableError(envName string) error {
	return fmt.Errorf("no value found for required environment variable: %s", envName)
}

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
	cmd.AddCommand(Generate())

	return cmd
}
