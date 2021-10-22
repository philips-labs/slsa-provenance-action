package github_test

import (
	"context"
	"os"
	"path"
	"runtime"
	"testing"

	gh "github.com/google/go-github/v39/github"
	"github.com/stretchr/testify/assert"

	"github.com/philips-labs/slsa-provenance-action/lib/github"
)

const (
	owner = "philips-labs"
	repo  = "slsa-provenance-action"
)

var githubToken string

func tokenRetriever() string {
	return os.Getenv("GITHUB_TOKEN")
}

func stringPointer(s string) *string {
	return &s
}

func boolPointer(b bool) *bool {
	return &b
}

func init() {
	githubToken = tokenRetriever()
}

func TestFetchRelease(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	client := createProvenanceClient(ctx)
	release, err := client.FetchRelease(ctx, owner, repo, "v0.1.1")

	if !assert.NoError(err) && assert.Nil(release) {
		return
	}
	assert.Equal(int64(51517953), release.GetID())
	assert.Equal("v0.1.1", release.GetTagName())
	assert.Len(release.Assets, 7)
}

func TestDownloadReleaseAssets(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()

	client := createProvenanceClient(ctx)

	release, err := client.FetchRelease(ctx, owner, repo, "v0.1.1")
	if !assert.NoError(err) && assert.Nil(release) {
		return
	}
	assert.Equal(int64(51517953), release.GetID())

	assets, err := client.DownloadReleaseAssets(ctx, owner, repo, release.GetID())
	if !assert.NoError(err) {
		return
	}
	defer func() {
		for _, a := range assets {
			if a.Content != nil {
				_ = a.Content.Close()
			}
		}
	}()

	assert.Len(assets, 7)
	assert.Equal("checksums.txt", assets[0].GetName())
	assert.NotNil(assets[0].Content)
	assert.Equal("slsa-provenance_0.1.1_linux_amd64.tar.gz", assets[1].GetName())
	assert.NotNil(assets[1].Content)
	assert.Equal("slsa-provenance_0.1.1_linux_arm64.tar.gz", assets[2].GetName())
	assert.NotNil(assets[2].Content)
	assert.Equal("slsa-provenance_0.1.1_macOS_amd64.tar.gz", assets[3].GetName())
	assert.NotNil(assets[3].Content)
	assert.Equal("slsa-provenance_0.1.1_macOS_arm64.tar.gz", assets[4].GetName())
	assert.NotNil(assets[4].Content)
	assert.Equal("slsa-provenance_0.1.1_windows_amd64.zip", assets[5].GetName())
	assert.NotNil(assets[5].Content)
	assert.Equal("slsa-provenance_0.1.1_windows_arm64.zip", assets[6].GetName())
	assert.NotNil(assets[6].Content)
}

func TestAddProvenanceToRelease(t *testing.T) {
	assert := assert.New(t)

	if githubToken == "" {
		t.Skip("skipping as GITHUB_TOKEN environment variable isn't set")
	}

	_, filename, _, _ := runtime.Caller(0)
	rootDir := path.Join(path.Dir(filename), "../..")
	provenanceFile := path.Join(rootDir, ".github/test_resource/example_build.provenance")

	ctx := context.Background()
	tc := github.NewOAuth2Client(ctx, tokenRetriever)
	client := github.NewProvenanceClient(tc)

	releaseId, err := createGitHubRelease(ctx, client, owner, repo, "v0.0.0-test")
	if !assert.NoError(err) {
		return
	}
	defer func() {
		_, err := client.Repositories.DeleteRelease(ctx, owner, repo, releaseId)
		assert.NoError(err)
	}()

	provenance, err := os.Open(provenanceFile)
	if !assert.NoError(err) && assert.Nil(provenance) {
		return
	}

	stat, err := provenance.Stat()
	if !assert.NoError(err) && assert.Nil(stat) {
		return
	}
	assert.Equal("example_build.provenance", stat.Name())

	asset, err := client.AddProvenanceToRelease(ctx, owner, repo, releaseId, provenance)
	if !assert.NoError(err) && assert.Nil(asset) {
		return
	}
	assert.Equal(stat.Name(), asset.GetName())
	assert.Equal("application/json; charset=utf-8", asset.GetContentType())
}

func TestListReleaseAssets(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	client := createProvenanceClient(ctx)
	opt := gh.ListOptions{PerPage: 2}
	assets, err := client.ListReleaseAssets(ctx, owner, repo, 51517953, opt)
	if !assert.NoError(err) {
		return
	}
	assert.Len(assets, 7)

	opt = gh.ListOptions{PerPage: 10}
	assets, err = client.ListReleaseAssets(ctx, owner, repo, 51517953, opt)
	if !assert.NoError(err) {
		return
	}
	assert.Len(assets, 7)
}

func TestListReleases(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	client := createProvenanceClient(ctx)
	opt := gh.ListOptions{PerPage: 1}
	releases, err := client.ListReleases(ctx, owner, repo, opt)
	if !assert.NoError(err) {
		return
	}
	assert.GreaterOrEqual(len(releases), 2)

	opt = gh.ListOptions{PerPage: 2}
	releases, err = client.ListReleases(ctx, owner, repo, opt)
	if !assert.NoError(err) {
		return
	}
	assert.GreaterOrEqual(len(releases), 2)
}

func createProvenanceClient(ctx context.Context) *github.ProvenanceClient {
	var client *github.ProvenanceClient
	if githubToken != "" {
		tc := github.NewOAuth2Client(ctx, tokenRetriever)
		client = github.NewProvenanceClient(tc)
	} else {
		client = github.NewProvenanceClient(nil)
	}
	return client
}

func createGitHubRelease(ctx context.Context, client *github.ProvenanceClient, owner, repo, version string, assets ...string) (int64, error) {
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
		client.AddProvenanceToRelease(ctx, owner, repo, rel.GetID(), asset)
	}

	return rel.GetID(), nil
}
