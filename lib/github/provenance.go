package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

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

// ReleaseEnvironment implements intoto.Provenancer to Generate provenance based on a GitHub release
type ReleaseEnvironment struct {
	*Environment
	rc      *ReleaseClient
	tagName string
}

// NewReleaseEnvironment creates a new instance of ReleaseEnvironment with the given tagName and provenanceClient
func NewReleaseEnvironment(gh Context, runner RunnerContext, tagName string, rc *ReleaseClient) *ReleaseEnvironment {
	return &ReleaseEnvironment{
		Environment: &Environment{
			Context: &gh,
			Runner:  &runner,
		},
		rc:      rc,
		tagName: tagName,
	}
}

// GenerateProvenanceStatement generates provenance from the GitHub release environment
//
// Release assets will be downloaded to the given artifactPath
//
// The artifactPath has to be a directory.
func (e *ReleaseEnvironment) GenerateProvenanceStatement(ctx context.Context, artifactPath string) (*intoto.Statement, error) {
	err := os.MkdirAll(artifactPath, 0755)
	if err != nil {
		return nil, err
	}
	isDir, err := isEmptyDirectory(artifactPath)
	if err != nil {
		return nil, err
	}
	if !isDir {
		return nil, errors.New("artifactPath has to be an empty directory")
	}

	owner := e.Context.RepositoryOwner
	repo := repositoryName(e.Context.Repository)
	rel, err := e.rc.FetchRelease(ctx, owner, repo, e.tagName)
	if err != nil {
		return nil, err
	}
	_, err = e.rc.DownloadReleaseAssets(ctx, owner, repo, rel.GetID(), artifactPath)
	if err != nil {
		return nil, err
	}

	return e.Environment.GenerateProvenanceStatement(ctx, artifactPath)
}

func isEmptyDirectory(p string) (bool, error) {
	f, err := os.Open(p)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

func repositoryName(repo string) string {
	repoParts := strings.Split(repo, "/")
	return repoParts[len(repoParts)-1]
}
