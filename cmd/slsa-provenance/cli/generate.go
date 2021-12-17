package cli

import (
	"github.com/spf13/cobra"
)

// Generate creates an instance of *cobra.Command to generate provenance
func Generate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate provenance using subcommands",
	}

	cmd.AddCommand(
		Files(),
		GitHubRelease(),
		OCI(),
	)

	return cmd
}
