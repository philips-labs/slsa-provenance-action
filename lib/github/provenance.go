package github

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/pkg/errors"

	"github.com/philips-labs/slsa-provenance-action/lib/intoto"
)

// GenerateProvenanceStatement generates provenance from the provided artifactPath
//
// The artifactPath can be a file or a directory.
func (e *Environment) GenerateProvenanceStatement(ctx context.Context, artifactPath string) (*intoto.Statement, error) {
	subjects, err := intoto.Subjects(artifactPath)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("resource path not found: [provided=%s]", artifactPath)
	} else if err != nil {
		return nil, err
	}

	repoURI := "https://github.com/" + e.Context.Repository

	event := AnyEvent{}
	if err := json.Unmarshal(e.Context.Event, &event); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal github context event json")
	}

	stmt := intoto.SLSAProvenanceStatement(
		intoto.WithSubject(subjects),
		intoto.WithBuilder(builderID(repoURI)),
		// NOTE: Re-runs are not uniquely identified and can cause run ID collisions.
		intoto.WithMetadata(fmt.Sprintf("%s/actions/runs/%s", repoURI, e.Context.RunID)),
		// NOTE: This is inexact as multiple workflows in a repo can have the same name.
		// See https://github.com/github/feedback/discussions/4188
		intoto.WithRecipe(
			RecipeType,
			e.Context.Workflow,
			nil,
			event.Inputs,
			[]intoto.Item{
				{URI: "git+" + repoURI, Digest: intoto.DigestSet{"sha1": e.Context.SHA}},
			},
		))

	return stmt, nil
}
