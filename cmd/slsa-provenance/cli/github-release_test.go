package cli_test

import (
	"context"
	"encoding/base64"
	"os"
	"path"
	"runtime"
	"testing"

	gh "github.com/google/go-github/v41/github"
	"github.com/stretchr/testify/assert"

	"github.com/philips-labs/slsa-provenance-action/cmd/slsa-provenance/cli"
	"github.com/philips-labs/slsa-provenance-action/lib/github"
)

func TestProvenenaceGitHubRelease(t *testing.T) {
	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		t.Skip("skipping as GITHUB_TOKEN environment variable isn't set")
	}
	assert := assert.New(t)

	_, filename, _, _ := runtime.Caller(0)
	rootDir := path.Join(path.Dir(filename), "../../..")
	artifactPath := path.Join(rootDir, "gh-release-test")
	provenanceFile := path.Join(artifactPath, "unittest.provenance")

	ctx := context.Background()
	owner, repo := "philips-labs", "slsa-provenance-action"
	oauthClient := github.NewOAuth2Client(ctx, func() string { return githubToken })
	client := github.NewReleaseClient(oauthClient)

	releaseID, err := createGitHubRelease(
		ctx,
		client,
		owner,
		repo,
		"v0.0.0-generate-test",
		path.Join(rootDir, "bin", "slsa-provenance"),
		path.Join(rootDir, "README.md"),
	)
	assert.NoError(err)

	defer func() {
		_ = os.RemoveAll(artifactPath)
		_, err = client.Repositories.DeleteRelease(ctx, owner, repo, releaseID)
	}()

	base64GitHubContext := base64.StdEncoding.EncodeToString([]byte(githubContext))
	base64RunnerContext := base64.StdEncoding.EncodeToString([]byte(runnerContext))

	output, err := executeCommand(cli.GitHubRelease(),
		"--artifact-path",
		artifactPath,
		"--github-context",
		base64GitHubContext,
		"--output-path",
		provenanceFile,
		"--runner-context",
		base64RunnerContext,
		"--tag-name",
		"v0.0.0-generate-test",
	)

	assert.NoError(err)
	assert.Contains(output, "Saving provenance to")

	if assert.FileExists(provenanceFile) {
		content, err := os.ReadFile(provenanceFile)
		assert.NoError(err)
		assert.Greater(len(content), 1)
	}
}

func createGitHubRelease(ctx context.Context, client *github.ReleaseClient, owner, repo, version string, assets ...string) (int64, error) {
	rel, _, err := client.Repositories.CreateRelease(
		ctx,
		owner,
		repo,
		&gh.RepositoryRelease{TagName: stringPointer(version), Name: stringPointer(version), Draft: boolPointer(true), Prerelease: boolPointer(true)},
	)
	if err != nil {
		return 0, err
	}

	for _, a := range assets {
		asset, err := os.Open(a)
		if err != nil {
			return 0, err
		}
		defer asset.Close()
		client.AddProvenanceToRelease(ctx, owner, repo, rel.GetID(), asset)
	}

	return rel.GetID(), nil
}
