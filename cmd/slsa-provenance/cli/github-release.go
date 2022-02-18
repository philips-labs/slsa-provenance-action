package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/philips-labs/slsa-provenance-action/cmd/slsa-provenance/cli/options"
	"github.com/philips-labs/slsa-provenance-action/internal/transport"
	"github.com/philips-labs/slsa-provenance-action/pkg/github"
	"github.com/philips-labs/slsa-provenance-action/pkg/intoto"
)

// GitHubRelease creates an instance of *cobra.Command to manage GitHub release provenance
func GitHubRelease() *cobra.Command {
	o := options.GitHubReleaseOptions{}

	cmd := &cobra.Command{
		Use:   "github-release",
		Short: "Generate provenance on GitHub release assets",
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

			tagName, err := o.GetTagName()
			if err != nil {
				return err
			}

			ghToken := os.Getenv("GITHUB_TOKEN")
			if ghToken == "" {
				return errors.New("GITHUB_TOKEN environment variable not set")
			}
			tc := github.NewOAuth2Client(cmd.Context(), func() string { return ghToken })
			tc.Transport = transport.TeeRoundTripper{
				RoundTripper: tc.Transport,
				Writer:       cmd.OutOrStdout(),
			}
			rc := github.NewReleaseClient(tc)
			env := github.NewReleaseEnvironment(*gh, *runner, tagName, rc, artifactPath)

			subjecter := intoto.NewFilePathSubjecter(artifactPath)
			stmt, err := env.GenerateProvenanceStatement(cmd.Context(), subjecter, materials...)
			if err != nil {
				return fmt.Errorf("failed to generate provenance: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Saving provenance to %s\n", outputPath)

			return env.PersistProvenanceStatement(cmd.Context(), stmt, outputPath)
		},
	}

	o.AddFlags(cmd)

	return cmd
}
