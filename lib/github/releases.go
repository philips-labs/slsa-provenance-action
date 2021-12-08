package github

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/google/go-github/v41/github"
	"golang.org/x/oauth2"
)

// TokenRetriever allows to implement a function to retrieve the token
// The token is placed in a StaticTokenSource to authenticate using oauth2.
type TokenRetriever func() string

// NewOAuth2Client creates a oauth2 client using the token from the TokenRetriever
func NewOAuth2Client(ctx context.Context, tokenRetriever TokenRetriever) *http.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: tokenRetriever()})
	return oauth2.NewClient(ctx, ts)
}

// ReleaseClient GitHub client adding convenience methods to add provenance to a release
type ReleaseClient struct {
	*github.Client
	httpClient *http.Client
}

// NewReleaseClient create new ReleaseClient instance
func NewReleaseClient(httpClient *http.Client) *ReleaseClient {
	return &ReleaseClient{
		Client:     github.NewClient(httpClient),
		httpClient: httpClient,
	}
}

// DownloadReleaseAssets download the assets for a release at the given storage location.
func (p *ReleaseClient) DownloadReleaseAssets(ctx context.Context, owner, repo, tag string, storageLocation string) ([]*github.ReleaseAsset, error) {
	downloadCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	relCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	r, _, err := p.Repositories.GetReleaseByTag(relCtx, owner, repo, tag)
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(storageLocation, 0755)
	if err != nil {
		return nil, err
	}

	assets := make([]*github.ReleaseAsset, len(r.Assets))

	for i, releaseAsset := range r.Assets {
		asset, _, err := p.Repositories.DownloadReleaseAsset(downloadCtx, owner, repo, releaseAsset.GetID(), p.httpClient)
		if err != nil {
			var errResponse *github.ErrorResponse
			if errors.As(err, &errResponse) {
				b, err := ioutil.ReadAll(errResponse.Response.Body)
				if err != nil {
					fmt.Println(b)
				}
			}
			return nil, err
		}
		err = saveFile(path.Join(storageLocation, releaseAsset.GetName()), asset)
		if err != nil {
			return nil, err
		}
		assets[i] = releaseAsset
	}

	return assets, nil
}

func saveFile(path string, content io.ReadCloser) error {
	assetFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer assetFile.Close()
	defer content.Close()

	_, err = io.Copy(assetFile, content)

	return err
}

// AddProvenanceToRelease uploads the provenance for the given release
func (p *ReleaseClient) AddProvenanceToRelease(ctx context.Context, owner, repo string, releaseID int64, provenance *os.File) (*github.ReleaseAsset, error) {
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
func (p *ReleaseClient) ListReleaseAssets(ctx context.Context, owner, repo string, releaseID int64, listOptions github.ListOptions) ([]*github.ReleaseAsset, error) {
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
func (p *ReleaseClient) ListReleases(ctx context.Context, owner, repo string, listOptions github.ListOptions) ([]*github.RepositoryRelease, error) {
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
