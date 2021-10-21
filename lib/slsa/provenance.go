package slsa

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/pkg/errors"

	"github.com/philips-labs/slsa-provenance-action/lib/github"
	"github.com/philips-labs/slsa-provenance-action/lib/intoto"
)

const (
	// GitHubHostedIDSuffix the GitHub hosted attestation type
	GitHubHostedIDSuffix = "/Attestations/GitHubHostedActions@v1"
	// SelfHostedIDSuffix the GitHub self hosted attestation type
	SelfHostedIDSuffix = "/Attestations/SelfHostedActions@v1"
	// RecipeType the attestion type for a recipe
	RecipeType = "https://github.com/Attestations/GitHubActionsWorkflow@v1"
	// PayloadContentType used to define the Envelope content type
	// See: https://github.com/in-toto/attestation#provenance-example
	PayloadContentType = "application/vnd.in-toto+json"
)

func builderID(repoURI string) string {
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		return repoURI + GitHubHostedIDSuffix
	}
	return repoURI + SelfHostedIDSuffix
}

// GenerateProvenanceStatement generates a in-toto provenance statement based on the github context
func GenerateProvenanceStatement(ctx context.Context, gh github.Context, runner github.RunnerContext, artifactPath string) (*intoto.Statement, error) {
	subjects, err := subjects(artifactPath)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("resource path not found: [provided=%s]", artifactPath)
	} else if err != nil {
		return nil, err
	}

	repoURI := "https://github.com/" + gh.Repository

	event := github.AnyEvent{}
	if err := json.Unmarshal(gh.Event, &event); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal github context event json")
	}

	stmt := intoto.SLSAProvenanceStatement(
		intoto.WithSubject(subjects),
		intoto.WithBuilder(builderID(repoURI)),
		// NOTE: Re-runs are not uniquely identified and can cause run ID collisions.
		intoto.WithMetadata(fmt.Sprintf("%s/actions/runs/%s", repoURI, gh.RunID)),
		// NOTE: This is inexact as multiple workflows in a repo can have the same name.
		// See https://github.com/github/feedback/discussions/4188
		intoto.WithRecipe(
			RecipeType,
			gh.Workflow,
			nil,
			event.Inputs,
			[]intoto.Item{
				{URI: "git+" + repoURI, Digest: intoto.DigestSet{"sha1": gh.SHA}},
			},
		))
	return stmt, nil
}
