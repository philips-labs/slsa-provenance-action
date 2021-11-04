package cli

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/peterbourgon/ff/v3/ffcli"
	"github.com/pkg/errors"

	"github.com/philips-labs/slsa-provenance-action/lib/github"
	"github.com/philips-labs/slsa-provenance-action/lib/intoto"
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
		outputPath     = flagset.String("output_path", "build.provenance", "The path to which the generated provenance should be written.")
		githubContext  = flagset.String("github_context", "", "The '${github}' context value.")
		runnerContext  = flagset.String("runner_context", "", "The '${runner}' context value.")
		extraMaterials = []string{}
	)
	flagset.Func("extra_materials", "Files that contain JSON encoded provenance to be included into the provenance", func(s string) error {
		extraMaterials = append(extraMaterials, strings.Fields(s)...)
		return nil
	})

	flagset.SetOutput(w)

	return &ffcli.Command{
		Name:       "generate",
		ShortUsage: "slsa-provenance generate",
		ShortHelp:  "Generates the slsa provenance file",
		FlagSet:    flagset,
		Exec: func(ctx context.Context, args []string) error {
			if *artifactPath == "" {
				flagset.Usage()
				return RequiredFlagError("-artifact_path")
			}
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
				return errors.Wrap(err, "failed to unmarshal github context json")
			}

			var runner github.RunnerContext
			if err := json.Unmarshal([]byte(*runnerContext), &runner); err != nil {
				return errors.Wrap(err, "failed to unmarshal runner context json")
			}

			ghToken := os.Getenv("GITHUB_TOKEN")
			if ghToken == "" {
				return errors.New("GITHUB_TOKEN environment variable not set")
			}

			tc := github.NewOAuth2Client(ctx, func() string { return ghToken })
			rc := github.NewReleaseClient(tc)
			env := createEnvironment(gh, runner, *tagName, rc)
			stmt, err := env.GenerateProvenanceStatement(ctx, *artifactPath)
			if err != nil {
				return errors.Wrap(err, "failed to generate provenance")
			}

			for _, extra := range extraMaterials {
				content, err := ioutil.ReadFile(extra)
				if err != nil {
					return errors.Wrapf(err, "Could not load extra materials from %s", extra)
				}
				var materials []intoto.Item
				if err = json.Unmarshal(content, &materials); err != nil {
					return errors.Wrapf(err, "Invalid JSON in extra materials file %s", extra)
				}
				stmt.Predicate.Materials = append(stmt.Predicate.Materials, materials...)
			}

			fmt.Fprintf(w, "Saving provenance to %s\n", *outputPath)

			return env.PersistProvenanceStatement(ctx, stmt, *outputPath)
		},
	}
}

func createEnvironment(gh github.Context, runner github.RunnerContext, tagName string, rc *github.ReleaseClient) intoto.Provenancer {
	if tagName != "" {
		return github.NewReleaseEnvironment(gh, runner, tagName, rc)
	}

	return &github.Environment{
		Context: &gh,
		Runner:  &runner,
	}
}
