package slsa

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"

	"github.com/philips-labs/slsa-provenance-action/lib/github"
	"github.com/philips-labs/slsa-provenance-action/lib/intoto"
)

const (
	gitHubHostedIDSuffix = "/Attestations/GitHubHostedActions@v1"
	selfHostedIDSuffix   = "/Attestations/SelfHostedActions@v1"
	recipeType           = "https://github.com/Attestations/GitHubActionsWorkflow@v1"
	payloadContentType   = "application/vnd.in-toto+json"
)

func builderID(repoURI string) string {
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		return repoURI + gitHubHostedIDSuffix
	}
	return repoURI + selfHostedIDSuffix
}

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
			recipeType,
			gh.Workflow,
			nil,
			event.Inputs,
			[]intoto.Item{
				{URI: "git+" + repoURI, Digest: intoto.DigestSet{"sha1": gh.SHA}},
			},
		))
	return stmt, nil
}
