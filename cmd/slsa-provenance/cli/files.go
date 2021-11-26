package cli

import (
	"flag"
	"fmt"
	"strings"

	"github.com/philips-labs/slsa-provenance-action/cmd/slsa-provenance/cli/options"
	"github.com/philips-labs/slsa-provenance-action/lib/github"
	"github.com/spf13/cobra"
)

// Files creates an instance of *ffcli.Command to manage file provenance
func Files() *cobra.Command {
	o := &options.FilesOptions{}

	var (
		flagset        = flag.NewFlagSet("slsa-provenance generate files", flag.ExitOnError)
		extraMaterials = []string{}
	)
	flagset.Func("extra_materials", "paths to files containing SLSA v0.1 formatted materials (JSON array) in to include in the provenance", func(s string) error {
		extraMaterials = append(extraMaterials, strings.Fields(s)...)
		return nil
	})

	cmd := &cobra.Command{
		Use:   "files",
		Short: fmt.Sprintf("%s generate files", cliName),
		Long:  "Generates slsa provenance for file(s)",
		RunE: func(cmd *cobra.Command, args []string) error {
			artifactPath, err := o.GetArtifactPath()
			if err != nil {
				return err
			}
			outputPath, err := o.GetOutputPath()
			if err != nil {
				return err
			}

			gh, err := o.GetGitHubContext()
			if err != nil {
				return err
			}

			runner, err := o.GetRunnerContext()
			if err != nil {
				return err
			}

			materials, err := o.GetExtraMaterials()
			if err != nil {
				return err
			}

			env := &github.Environment{
				Context: gh,
				Runner:  runner,
			}

			stmt, err := env.GenerateProvenanceStatement(cmd.Context(), artifactPath)
			if err != nil {
				return fmt.Errorf("failed to generate provenance: %w", err)
			}

			stmt.Predicate.Materials = append(stmt.Predicate.Materials, materials...)

			fmt.Fprintf(cmd.OutOrStdout(), "Saving provenance to %s\n", outputPath)

			return env.PersistProvenanceStatement(cmd.Context(), stmt, outputPath)
		},
	}

	o.AddFlags(cmd)

	return cmd
}
