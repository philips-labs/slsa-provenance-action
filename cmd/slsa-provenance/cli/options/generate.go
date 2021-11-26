package options

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/philips-labs/slsa-provenance-action/lib/github"
	"github.com/philips-labs/slsa-provenance-action/lib/intoto"
)

type GenerateOptions struct {
	GitHubContext  string
	RunnerContext  string
	ExtraMaterials []string
}

func (o *GenerateOptions) GetGitHubContext() (*github.Context, error) {
	if o.GitHubContext == "" {
		return nil, RequiredFlagError("github-context")
	}
	var gh github.Context
	if err := json.Unmarshal([]byte(o.GitHubContext), &gh); err != nil {
		return nil, fmt.Errorf("failed to unmarshal github context json: %w", err)
	}
	return &gh, nil
}

func (o *GenerateOptions) GetRunnerContext() (*github.RunnerContext, error) {
	if o.RunnerContext == "" {
		return nil, RequiredFlagError("runner-context")
	}
	var runner github.RunnerContext
	if err := json.Unmarshal([]byte(o.RunnerContext), &runner); err != nil {
		return nil, fmt.Errorf("failed to unmarshal runner context json: %w", err)
	}
	return &runner, nil
}

func (o *GenerateOptions) GetExtraMaterials() ([]intoto.Item, error) {
	var materials []intoto.Item

	for _, extra := range o.ExtraMaterials {
		content, err := os.ReadFile(extra)
		if err != nil {
			return nil, fmt.Errorf("failed retrieving extra materials: %w", err)
		}
		if err = json.Unmarshal(content, &materials); err != nil {
			return nil, fmt.Errorf("failed retrieving extra materials: invalid JSON in extra materials file %s: %w", extra, err)
		}
		for _, material := range materials {
			if material.URI == "" {
				return nil, fmt.Errorf("failed retrieving extra materials: empty or missing \"uri\" field in %s", extra)
			}
			if len(material.Digest) == 0 {
				return nil, fmt.Errorf("failed retrieving extra materials: empty or missing \"digest\" in %s", extra)
			}
			materials = append(materials, material)
		}
	}

	return materials, nil
}

func (o *GenerateOptions) AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&o.GitHubContext, "github-context", "", "The '${github}' context value.")
	cmd.PersistentFlags().StringVar(&o.RunnerContext, "runner-context", "", "The '${runner}' context value.")
	cmd.PersistentFlags().StringSliceVarP(&o.ExtraMaterials, "extra-materials", "m", nil, "The '${runner}' context value.")
}
