package github

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
)

// TokenRetriever allows to implement a function to retrieve the token
// The token is placed in a StaticTokenSource to authenticate using oauth2.
type TokenRetriever func() string

func NewOAuth2Client(ctx context.Context, tokenRetriever TokenRetriever) *http.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: tokenRetriever()})
	return oauth2.NewClient(ctx, ts)
}

// ProvenanceClient GitHub client adding convenience methods to add provenance to a release
type ProvenanceClient struct {
	*github.Client
}

// NewProvenanceClient create new ProvenanceClient instance
func NewProvenanceClient(httpClient *http.Client) *ProvenanceClient {
	return &ProvenanceClient{
		Client: github.NewClient(httpClient),
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
