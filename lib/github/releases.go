package github

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
)

// ReleaseAsset holds the release asset information and it's contents.
type ReleaseAsset struct {
	*github.ReleaseAsset
	Content io.ReadCloser
}

// TokenRetriever allows to implement a function to retrieve the token
// The token is placed in a StaticTokenSource to authenticate using oauth2.
type TokenRetriever func() string

// NewOAuth2Client creates a oauth2 client using the token from the TokenRetriever
func NewOAuth2Client(ctx context.Context, tokenRetriever TokenRetriever) *http.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: tokenRetriever()})
	return oauth2.NewClient(ctx, ts)
}

// ProvenanceClient GitHub client adding convenience methods to add provenance to a release
type ProvenanceClient struct {
	*github.Client
	httpClient *http.Client
}

// NewProvenanceClient create new ProvenanceClient instance
func NewProvenanceClient(httpClient *http.Client) *ProvenanceClient {
	return &ProvenanceClient{
		Client:     github.NewClient(httpClient),
		httpClient: httpClient,
	}
}

// FetchRelease get the release by its tagName
func (p *ProvenanceClient) FetchRelease(ctx context.Context, owner, repo, tagName string) (*github.RepositoryRelease, error) {
	client := p.Client

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// TODO: add pagination when there are tons of releases for the repo
	releases, _, err := client.Repositories.ListReleases(ctx, owner, repo, &github.ListOptions{})
	if err != nil {
		return nil, err
	}

	var rel *github.RepositoryRelease
	for _, r := range releases {
		if *r.TagName == tagName {
			rel = r
			break
		}
	}

	return rel, nil
}

// DownloadReleaseAssets download the assets for a release.
// It is up to the caller to Close the ReadCloser.
func (p *ProvenanceClient) DownloadReleaseAssets(ctx context.Context, owner, repo string, releaseID int64) ([]ReleaseAsset, error) {
	// TODO: add pagination when there are tons of releaseAssets not fitting in a single page for the release
	releaseAssets, _, err := p.Repositories.ListReleaseAssets(ctx, owner, repo, releaseID, &github.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list release assets: %w", err)
	}
	assets := make([]ReleaseAsset, len(releaseAssets))

	for i, releaseAsset := range releaseAssets {
		asset, _, err := p.Repositories.DownloadReleaseAsset(ctx, owner, repo, releaseAsset.GetID(), p.httpClient)
		if err != nil {
			return nil, err
		}
		assets[i] = ReleaseAsset{
			ReleaseAsset: releaseAsset,
			Content:      asset,
		}
	}

	return assets, nil
}

// AddProvenanceToRelease uploads the provenance for the given release
func (p *ProvenanceClient) AddProvenanceToRelease(ctx context.Context, owner, repo string, releaseID int64, provenance *os.File) (*github.ReleaseAsset, error) {
	client := p.Client

	stat, err := provenance.Stat()
	if err != nil {
		return nil, err
	}
	uploadOptions := &github.UploadOptions{Name: stat.Name(), MediaType: "application/json; charset=utf-8"}
	asset, _, err := client.Repositories.UploadReleaseAsset(ctx, owner, repo, releaseID, uploadOptions, provenance)
	return asset, err
}
