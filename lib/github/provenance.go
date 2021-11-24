package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

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
		return nil, fmt.Errorf("failed to unmarshal github context event json: %w", err)
	}

	stmt := intoto.SLSAProvenanceStatement(
		intoto.WithSubject(subjects),
		intoto.WithBuilder(builderID(repoURI)),
		// NOTE: Re-runs are not uniquely identified and can cause run ID collisions.
		intoto.WithMetadata(fmt.Sprintf("%s/actions/runs/%s", repoURI, e.Context.RunID)),
		// NOTE: This is inexact as multiple workflows in a repo can have the same name.
		// See https://github.com/github/feedback/discussions/4188
		intoto.WithInvocation(
			BuildType,
			e.Context.Workflow,
			nil,
			event.Inputs,
			[]intoto.Item{
				{URI: "git+" + repoURI, Digest: intoto.DigestSet{"sha1": e.Context.SHA}},
			},
		))

	return stmt, nil
}

// PersistProvenanceStatement writes the provenance statement at the given path
func (e *Environment) PersistProvenanceStatement(ctx context.Context, stmt *intoto.Statement, path string) error {
	// NOTE: At L1, writing the in-toto Statement type is sufficient but, at
	// higher SLSA levels, the Statement must be encoded and wrapped in an
	// Envelope to support attaching signatures.
	payload, err := json.MarshalIndent(stmt, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal provenance: %w", err)
	}
	if err := os.WriteFile(path, payload, 0755); err != nil {
		return fmt.Errorf("failed to write provenance: %w", err)
	}

	return nil
}

// ReleaseEnvironment implements intoto.Provenancer to Generate provenance based on a GitHub release
type ReleaseEnvironment struct {
	*Environment
	rc        *ReleaseClient
	tagName   string
	releaseID int64
}

// NewReleaseEnvironment creates a new instance of ReleaseEnvironment with the given tagName and provenanceClient
func NewReleaseEnvironment(gh Context, runner RunnerContext, tagName string, rc *ReleaseClient) *ReleaseEnvironment {
	return &ReleaseEnvironment{
		Environment: &Environment{
			Context: &gh,
			Runner:  &runner,
		},
		rc:        rc,
		tagName:   tagName,
		releaseID: 0,
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

	releaseID, err := e.GetReleaseID(ctx, e.tagName)
	if err != nil {
		return nil, err
	}
	_, err = e.rc.DownloadReleaseAssets(ctx, owner, repo, releaseID, artifactPath)
	if err != nil {
		return nil, err
	}

	return e.Environment.GenerateProvenanceStatement(ctx, artifactPath)
}

// PersistProvenanceStatement writes the provenance statement at the given path and uploads it to the GitHub release
func (e *ReleaseEnvironment) PersistProvenanceStatement(ctx context.Context, stmt *intoto.Statement, path string) error {
	err := e.Environment.PersistProvenanceStatement(ctx, stmt, path)
	if err != nil {
		return err
	}

	stmtFile, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open provenance statement: %w", err)
	}

	owner := e.Context.RepositoryOwner
	repo := repositoryName(e.Context.Repository)
	_, err = e.rc.AddProvenanceToRelease(ctx, owner, repo, e.releaseID, stmtFile)
	if err != nil {
		return fmt.Errorf("failed to upload provenance to release: %w", err)
	}

	return nil
}

// GetReleaseID fetches a release and caches the releaseID in the environment
func (e *ReleaseEnvironment) GetReleaseID(ctx context.Context, tagName string) (int64, error) {
	owner := e.Context.RepositoryOwner
	repo := repositoryName(e.Context.Repository)

	if e.releaseID == 0 {
		rel, err := e.rc.FetchRelease(ctx, owner, repo, e.tagName)
		if err != nil {
			return 0, err
		}
		e.releaseID = rel.GetID()
	}

	return e.releaseID, nil
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
