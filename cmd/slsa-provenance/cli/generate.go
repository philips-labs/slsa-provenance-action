package cli

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/peterbourgon/ff/v3/ffcli"
	"github.com/pkg/errors"

	"github.com/philips-labs/slsa-provenance-action/lib/provenance"
)

const (
	gitHubHostedIDSuffix = "/Attestations/GitHubHostedActions@v1"
	selfHostedIDSuffix   = "/Attestations/SelfHostedActions@v1"
	typeID               = "https://github.com/Attestations/GitHubActionsWorkflow@v1"
	payloadContentType   = "application/vnd.in-toto+json"
)

// RequiredFlagError creates an error flag error for the given flag name
func RequiredFlagError(flagName string) error {
	return fmt.Errorf("no value found for required flag: %s", flagName)
}

// subjects walks the file or directory at "root" and hashes all files.
func subjects(root string) ([]provenance.Subject, error) {
	var s []provenance.Subject
	return s, filepath.Walk(root, func(abspath string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		relpath, err := filepath.Rel(root, abspath)
		if err != nil {
			return err
		}
		// Note: filepath.Rel() returns "." when "root" and "abspath" point to the same file.
		if relpath == "." {
			relpath = filepath.Base(root)
		}
		contents, err := ioutil.ReadFile(abspath)
		if err != nil {
			return err
		}
		sha := sha256.Sum256(contents)
		shaHex := hex.EncodeToString(sha[:])
		s = append(s, provenance.Subject{Name: relpath, Digest: provenance.DigestSet{"sha256": shaHex}})
		return nil
	})
}

// Generate creates an instance of *ffcli.Command to generate provenance
func Generate() *ffcli.Command {
	var (
		flagset       = flag.NewFlagSet("slsa-provenance version", flag.ExitOnError)
		artifactPath  = flagset.String("artifact_path", "", "The file or dir path of the artifacts for which provenance should be generated.")
		outputPath    = flagset.String("output_path", "build.provenance", "The path to which the generated provenance should be written.")
		githubContext = flagset.String("github_context", "", "The '${github}' context value.")
		runnerContext = flagset.String("runner_context", "", "The '${runner}' context value.")
	)

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

			stmt := provenance.Statement{PredicateType: "https://slsa.dev/provenance/v0.1", Type: "https://in-toto.io/Statement/v0.1"}
			subjects, err := subjects(*artifactPath)
			if os.IsNotExist(err) {
				return fmt.Errorf("resource path not found: [provided=%s]", *artifactPath)
			} else if err != nil {
				return err
			}

			stmt.Subject = append(stmt.Subject, subjects...)
			stmt.Predicate = provenance.Predicate{
				Builder: provenance.Builder{},
				Metadata: provenance.Metadata{
					Completeness: provenance.Completeness{
						Arguments:   true,
						Environment: false,
						Materials:   false,
					},
					Reproducible:    false,
					BuildFinishedOn: time.Now().UTC().Format(time.RFC3339),
				},
				Recipe: provenance.Recipe{
					Type:              typeID,
					DefinedInMaterial: 0,
				},
				Materials: []provenance.Item{},
			}

			context := provenance.AnyContext{}
			if err := json.Unmarshal([]byte(*githubContext), &context.GitHubContext); err != nil {
				return errors.Wrap(err, "failed to unmarshal github context json")
			}
			if err := json.Unmarshal([]byte(*runnerContext), &context.RunnerContext); err != nil {
				return errors.Wrap(err, "failed to unmarshal runner context json")
			}
			gh := context.GitHubContext
			// Remove access token from the generated provenance.
			context.GitHubContext.Token = ""
			// NOTE: Re-runs are not uniquely identified and can cause run ID collisions.
			repoURI := "https://github.com/" + gh.Repository
			stmt.Predicate.Metadata.BuildInvocationID = repoURI + "/actions/runs/" + gh.RunID
			// NOTE: This is inexact as multiple workflows in a repo can have the same name.
			// See https://github.com/github/feedback/discussions/4188
			stmt.Predicate.Recipe.EntryPoint = gh.Workflow
			event := provenance.AnyEvent{}
			if err := json.Unmarshal(context.GitHubContext.Event, &event); err != nil {
				return errors.Wrap(err, "failed to unmarshal github context event json")
			}

			stmt.Predicate.Recipe.Arguments = event.Inputs
			stmt.Predicate.Materials = append(stmt.Predicate.Materials, provenance.Item{URI: "git+" + repoURI, Digest: provenance.DigestSet{"sha1": gh.SHA}})

			if os.Getenv("GITHUB_ACTIONS") == "true" {
				stmt.Predicate.Builder.ID = repoURI + gitHubHostedIDSuffix
			} else {
				stmt.Predicate.Builder.ID = repoURI + selfHostedIDSuffix
			}

			// NOTE: At L1, writing the in-toto Statement type is sufficient but, at
			// higher SLSA levels, the Statement must be encoded and wrapped in an
			// Envelope to support attaching signatures.
			payload, _ := json.MarshalIndent(stmt, "", "  ")
			fmt.Println("Provenance:\n" + string(payload))
			if err := ioutil.WriteFile(*outputPath, payload, 0755); err != nil {
				return errors.Wrap(err, "failed to write provenance")
			}

			return nil
		},
	}
}
