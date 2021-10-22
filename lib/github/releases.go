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
	listCtx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	allReleases, err := p.ListReleases(listCtx, owner, repo, github.ListOptions{PerPage: 10})
	if err != nil {
		return nil, err
	}

	var rel *github.RepositoryRelease
	for _, r := range allReleases {
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
	listCtx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	allAssets, err := p.ListReleaseAssets(listCtx, owner, repo, releaseID, github.ListOptions{PerPage: 10})
	if err != nil {
		return nil, err
	}
	assets := make([]ReleaseAsset, len(allAssets))

	downloadCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()
	for i, releaseAsset := range allAssets {
		asset, _, err := p.Repositories.DownloadReleaseAsset(downloadCtx, owner, repo, releaseAsset.GetID(), p.httpClient)
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
	uploadCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()
	asset, _, err := client.Repositories.UploadReleaseAsset(uploadCtx, owner, repo, releaseID, uploadOptions, provenance)
	return asset, err
}

// ListReleaseAssets will retrieve the list of all release assets.
func (p *ProvenanceClient) ListReleaseAssets(ctx context.Context, owner, repo string, releaseID int64, listOptions github.ListOptions) ([]*github.ReleaseAsset, error) {
	var allAssets []*github.ReleaseAsset
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		assets, resp, err := p.Repositories.ListReleaseAssets(ctx, owner, repo, releaseID, &listOptions)
		if err != nil {
			return nil, fmt.Errorf("failed to list release assets: %w", err)
		}
		allAssets = append(allAssets, assets...)
		if resp.NextPage == 0 {
			break
		}
		listOptions.Page = resp.NextPage
	}
	return allAssets, nil
}

// ListReleases will retrieve the list of all releases.
func (p *ProvenanceClient) ListReleases(ctx context.Context, owner, repo string, listOptions github.ListOptions) ([]*github.RepositoryRelease, error) {
	var allReleases []*github.RepositoryRelease
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		releases, resp, err := p.Repositories.ListReleases(ctx, owner, repo, &listOptions)
		if err != nil {
			return nil, fmt.Errorf("failed to list releases: %w", err)
		}
		allReleases = append(allReleases, releases...)
		if resp.NextPage == 0 {
			break
		}
		listOptions.Page = resp.NextPage
	}
	return allReleases, nil
}
