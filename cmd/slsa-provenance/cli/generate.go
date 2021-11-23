package cli

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/peterbourgon/ff/v3/ffcli"

	"github.com/philips-labs/slsa-provenance-action/lib/github"
)

// RequiredFlagError creates a required flag error for the given flag name
func RequiredFlagError(flagName string) error {
	return fmt.Errorf("no value found for required flag: %s", flagName)
}

// Generate creates an instance of *ffcli.Command to generate provenance
func Generate(w io.Writer) *ffcli.Command {
	var (
		flagset = flag.NewFlagSet("slsa-provenance generate", flag.ExitOnError)
		tagName = flagset.String("tag_name", "", `The github release to generate provenance on.
(if set the artifacts will be downloaded from the release and the provenance will be added as an additional release asset.)`)
		artifactPath   = flagset.String("artifact_path", "", "The file or dir path of the artifacts for which provenance should be generated.")
		outputPath     = flagset.String("output_path", "provenance.json", "The path to which the generated provenance should be written.")
		githubContext  = flagset.String("github_context", "", "The '${github}' context value.")
		runnerContext  = flagset.String("runner_context", "", "The '${runner}' context value.")
		extraMaterials = []string{}
	)
	flagset.Func("extra_materials", "paths to files containing SLSA v0.1 formatted materials (JSON array) in to include in the provenance", func(s string) error {
		extraMaterials = append(extraMaterials, strings.Fields(s)...)
		return nil
	})

	flagset.SetOutput(w)

	return &ffcli.Command{
		Name:       "generate",
		ShortUsage: "slsa-provenance generate",
		ShortHelp:  "Generates the slsa provenance file",
		FlagSet:    flagset,
		Subcommands: []*ffcli.Command{
			Files(w),
			GitHubRelease(w),
		},
		Exec: func(ctx context.Context, args []string) error {
			if *outputPath == "" {
				flagset.Usage()
				return RequiredFlagError("-output_path")
			}
			if *githubContext == "" {
				flagset.Usage()
				return RequiredFlagError("-github_context")
			}
			if *runnerContext == "" {
				flagset.Usage()
				return RequiredFlagError("-runner_context")
			}

			var gh github.Context
			if err := json.Unmarshal([]byte(*githubContext), &gh); err != nil {
				return fmt.Errorf("failed to unmarshal github context json: %w", err)
			}

			var runner github.RunnerContext
			if err := json.Unmarshal([]byte(*runnerContext), &runner); err != nil {
				return fmt.Errorf("failed to unmarshal runner context json: %w", err)
			}

			return nil
		},
	}
}
