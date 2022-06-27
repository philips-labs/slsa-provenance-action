package options

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/philips-labs/slsa-provenance-action/pkg/github"
	"github.com/philips-labs/slsa-provenance-action/pkg/intoto"
)

// GenerateOptions Commandline flags used for the generate command.
type GenerateOptions struct {
	GitHubContext  string
	RunnerContext  string
	OutputPath     string
	ExtraMaterials []string
}

const (
	// ContextLen defines the context content limit.
	ContextLen = 1024 * 1024 // 1 MB
)

var (
	// EnvGithubContext holds the environment variable name for Github context.
	EnvGithubContext = "GITHUB_CONTEXT"
	// EnvRunnerContext holds the environment variable name for Runner context.
	EnvRunnerContext = "RUNNER_CONTEXT"
)

// GetGitHubContext The '${github}' context value, retrieved in a GitHub workflow.
func (o *GenerateOptions) GetGitHubContext() (*github.Context, error) {
	// Retrieve context by environment
	githubContext := os.Getenv(EnvGithubContext)
	if githubContext == "" {
		return nil, RequiredEnvironmentVariableError(EnvGithubContext)
	}

	// 1MB should be more than enough
	lr := io.LimitReader(strings.NewReader(githubContext), ContextLen)

	// Decode context
	var gh github.Context
	if err := json.NewDecoder(lr).Decode(&gh); err != nil {
		return nil, fmt.Errorf("failed to unmarshal github context json: %w", err)
	}

	// No error
	return &gh, nil
}

// GetRunnerContext The '${runner}' context value, retrieved in a GitHub workflow.
func (o *GenerateOptions) GetRunnerContext() (*github.RunnerContext, error) {
	// Retrieve context by environment
	runnerContext := os.Getenv(EnvRunnerContext)
	if runnerContext == "" {
		return nil, RequiredEnvironmentVariableError(EnvRunnerContext)
	}

	// 1MB should be more than enough
	lr := io.LimitReader(strings.NewReader(runnerContext), ContextLen)

	// Decode context
	var runner github.RunnerContext
	if err := json.NewDecoder(lr).Decode(&runner); err != nil {
		return nil, fmt.Errorf("failed to unmarshal runner context json: %w", err)
	}

	return &runner, nil
}

// GetOutputPath The location to write the provenance file.
func (o *GenerateOptions) GetOutputPath() (string, error) {
	if o.OutputPath == "" {
		return "", RequiredFlagError("output-path")
	}
	return o.OutputPath, nil
}

// GetExtraMaterials Additional material files to be used when generating provenance.
func (o *GenerateOptions) GetExtraMaterials() ([]intoto.Item, error) {
	var materials []intoto.Item

	for _, extra := range o.ExtraMaterials {
		file, err := os.Open(extra)
		if err != nil {
			return nil, fmt.Errorf("failed retrieving extra materials: %w", err)
		}
		defer file.Close()

		m, err := intoto.ReadMaterials(file)
		if err != nil {
			return nil, fmt.Errorf("failed retrieving extra materials for %s: %w", extra, err)
		}
		materials = append(materials, m...)
	}

	return materials, nil
}

// AddFlags Registers the flags with the cobra.Command.
func (o *GenerateOptions) AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&o.OutputPath, "output-path", "provenance.json", "The path to which the generated provenance should be written.")
	cmd.PersistentFlags().StringSliceVarP(&o.ExtraMaterials, "extra-materials", "m", nil, "The '${runner}' context value.")
}
