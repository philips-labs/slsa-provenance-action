package github_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path"
	"runtime"
	"strings"
	"testing"

	gh "github.com/google/go-github/v41/github"
	"github.com/stretchr/testify/assert"

	"github.com/philips-labs/slsa-provenance-action/internal/transport"
	"github.com/philips-labs/slsa-provenance-action/pkg/github"
)

const (
	owner = "philips-labs"
	repo  = "slsa-provenance-action"
)

const (
	releasesAPI       = "https://api.github.com/repos/philips-labs/slsa-provenance-action/releases"
	firstRelease      = "v0.1.1"
	assetsURLTemplate = "GET: %s/assets/%d"
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

	client, requestLogger := createReleaseClient(ctx)
	release, err := client.FetchRelease(ctx, owner, repo, firstRelease)

	if !assert.NoError(err) && !assert.NotNil(release) {
		return
	}
	assert.Equal(int64(51517953), release.GetID())
	assert.Equal(firstRelease, release.GetTagName())
	assert.Len(release.Assets, 7)
	assert.Equal(fmt.Sprintf("GET: %s?per_page=25\n", releasesAPI), requestLogger.String())
}

func TestDownloadReleaseAssets(t *testing.T) {
	if tokenRetriever() == "" {
		t.Skip("skipping as GITHUB_TOKEN environment variable isn't set")
	}
	assert := assert.New(t)

	ctx := context.Background()

	api := releasesAPI
	version := firstRelease
	client, requestLogger := createReleaseClient(ctx)
	release, _, err := client.Repositories.GetReleaseByTag(ctx, owner, repo, version)
	if !assert.NoError(err) || !assert.NotNil(release) {
		return
	}
	assert.Equal(int64(51517953), release.GetID())
	assert.NotEmpty(requestLogger)
	assert.Equal(fmt.Sprintf("GET: %s/tags/%s\n", api, version), requestLogger.String())
	requestLogger.Reset()

	_, filename, _, _ := runtime.Caller(0)
	rootDir := path.Join(path.Dir(filename), "../..")
	artifactPath := path.Join(rootDir, "download-test")
	assets, err := client.DownloadReleaseAssets(ctx, owner, repo, release.GetID(), artifactPath)
	if !assert.NoError(err) {
		return
	}
	assert.NotEmpty(requestLogger)
	assert.Contains(requestLogger.String(), fmt.Sprintf("GET: %s/%d/assets?per_page=10", api, release.GetID()))
	assert.Contains(requestLogger.String(), fmt.Sprintf(assetsURLTemplate, api, assets[0].GetID()))
	assert.Contains(requestLogger.String(), fmt.Sprintf(assetsURLTemplate, api, assets[1].GetID()))
	assert.Contains(requestLogger.String(), fmt.Sprintf(assetsURLTemplate, api, assets[2].GetID()))
	assert.Contains(requestLogger.String(), fmt.Sprintf(assetsURLTemplate, api, assets[3].GetID()))
	assert.Contains(requestLogger.String(), fmt.Sprintf(assetsURLTemplate, api, assets[4].GetID()))
	assert.Contains(requestLogger.String(), fmt.Sprintf(assetsURLTemplate, api, assets[5].GetID()))
	assert.Contains(requestLogger.String(), fmt.Sprintf(assetsURLTemplate, api, assets[6].GetID()))
	assert.Contains(requestLogger.String(), "GET: https://objects.githubusercontent.com/github-production-release-asset")

	defer func() {
		_ = os.RemoveAll(artifactPath)
	}()

	assert.Len(assets, 7)
	assert.Equal("checksums.txt", assets[0].GetName())
	assert.FileExists(path.Join(artifactPath, assets[0].GetName()))
	assert.Equal("slsa-provenance_0.1.1_linux_amd64.tar.gz", assets[1].GetName())
	assert.FileExists(path.Join(artifactPath, assets[1].GetName()))
	assert.Equal("slsa-provenance_0.1.1_linux_arm64.tar.gz", assets[2].GetName())
	assert.FileExists(path.Join(artifactPath, assets[2].GetName()))
	assert.Equal("slsa-provenance_0.1.1_macOS_amd64.tar.gz", assets[3].GetName())
	assert.FileExists(path.Join(artifactPath, assets[3].GetName()))
	assert.Equal("slsa-provenance_0.1.1_macOS_arm64.tar.gz", assets[4].GetName())
	assert.FileExists(path.Join(artifactPath, assets[4].GetName()))
	assert.Equal("slsa-provenance_0.1.1_windows_amd64.zip", assets[5].GetName())
	assert.FileExists(path.Join(artifactPath, assets[5].GetName()))
	assert.Equal("slsa-provenance_0.1.1_windows_arm64.zip", assets[6].GetName())
	assert.FileExists(path.Join(artifactPath, assets[6].GetName()))
}

func TestAddProvenanceToRelease(t *testing.T) {
	assert := assert.New(t)

	if githubToken == "" {
		t.Skip("skipping as GITHUB_TOKEN environment variable isn't set")
	}

	_, filename, _, _ := runtime.Caller(0)
	rootDir := path.Join(path.Dir(filename), "../..")
	provenanceFile := path.Join(rootDir, ".github/test_resource/example_provenance.json")

	ctx := context.Background()
	tc := github.NewOAuth2Client(ctx, tokenRetriever)
	client := github.NewReleaseClient(tc)

	releaseId, err := createGitHubRelease(ctx, client, owner, repo, "v0.0.0-test")
	if !assert.NoError(err) {
		return
	}
	defer func() {
		_, err := client.Repositories.DeleteRelease(ctx, owner, repo, releaseId)
		assert.NoError(err)
	}()

	provenance, err := os.Open(provenanceFile)
	if !assert.NoError(err) || !assert.NotNil(provenance) {
		return
	}
	defer provenance.Close()

	stat, err := provenance.Stat()
	if !assert.NoError(err) || !assert.NotNil(stat) {
		return
	}
	assert.Equal("example_provenance.json", stat.Name())

	asset, err := client.AddProvenanceToRelease(ctx, owner, repo, releaseId, provenance)
	if !assert.NoError(err) || !assert.NotNil(asset) {
		return
	}
	assert.Equal(stat.Name(), asset.GetName())
	assert.Equal("application/json; charset=utf-8", asset.GetContentType())
}

func TestListReleaseAssets(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	api := releasesAPI + "/51517953/assets"

	client, requestLogger := createReleaseClient(ctx)
	opt := gh.ListOptions{PerPage: 4}
	assets, err := client.ListReleaseAssets(ctx, owner, repo, 51517953, opt)
	if !assert.NoError(err) {
		return
	}
	assert.Len(assets, 7)
	assert.NotEmpty(requestLogger)
	assert.Equal(expectedRequestPages("GET", api, opt.PerPage, 2), requestLogger.String())

	opt = gh.ListOptions{PerPage: 10}
	assets, err = client.ListReleaseAssets(ctx, owner, repo, 51517953, opt)
	if !assert.NoError(err) {
		return
	}
	assert.Len(assets, 7)

	_, err = client.ListReleaseAssets(ctx, owner, repo, 0, opt)
	assert.EqualError(err, "failed to list release assets: GET "+releasesAPI+"/0/assets?per_page=10: 404 Not Found []")
}

func TestListReleases(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	api := releasesAPI

	client, requestLogger := createReleaseClient(ctx)
	opt := gh.ListOptions{PerPage: 25}
	releases, err := client.ListReleases(ctx, owner, repo, opt)
	if !assert.NoError(err) {
		return
	}
	assert.NotEmpty(requestLogger)
	assert.Equal(expectedRequestPages("GET", api, opt.PerPage, 1), requestLogger.String())
	assert.GreaterOrEqual(len(releases), 18)

	opt = gh.ListOptions{PerPage: 2}
	releases, err = client.ListReleases(ctx, owner, repo, opt)
	if !assert.NoError(err) {
		return
	}
	assert.GreaterOrEqual(len(releases), 18)

	opt = gh.ListOptions{PerPage: 2}
	_, err = client.ListReleases(ctx, owner, repo+"-fake", opt)
	assert.EqualError(err, "failed to list releases: GET https://api.github.com/repos/philips-labs/slsa-provenance-action-fake/releases?per_page=2: 404 Not Found []")
}

func expectedRequestPages(method string, api string, size int, count int) string {
	w := &strings.Builder{}

	for i := 1; i <= count; i++ {
		if i == 1 {
			fmt.Fprintf(w, "%s: %s?per_page=%d\n", method, api, size)
		} else {
			fmt.Fprintf(w, "%s: %s?page=%d&per_page=%d\n", method, api, i, size)
		}
	}

	return w.String()
}

func createReleaseClient(ctx context.Context) (*github.ReleaseClient, *strings.Builder) {
	var client *github.ReleaseClient
	var writer strings.Builder

	if githubToken != "" {
		tc := github.NewOAuth2Client(ctx, tokenRetriever)
		tc.Transport = transport.TeeRoundTripper{
			RoundTripper: tc.Transport,
			Writer:       &writer,
		}
		client = github.NewReleaseClient(tc)
	} else {
		client = github.NewReleaseClient(&http.Client{
			Transport: transport.TeeRoundTripper{
				RoundTripper: http.DefaultTransport,
				Writer:       &writer,
			},
		})
	}
	return client, &writer
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
