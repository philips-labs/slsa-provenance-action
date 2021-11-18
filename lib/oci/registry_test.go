package oci

import (
	"context"
	"testing"

	"github.com/docker/docker/client"
	"github.com/stretchr/testify/assert"
)

func TestPullImageTags(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if !assert.NoError(err) {
		return
	}
	repo := "ghcr.io/philips-labs/slsa-provenance"
	expectedDigest := "sha256:194b471a878add368bf02a7935fa099024576c029491bcefaeb87f81efa093a3"
	expectedTags := []string{"33ba3da2213c83ce02df0f2f6ba925ec79037f9d", "v0.4.0"}
	s := NewContainerSubjecter(cli, repo, expectedDigest, expectedTags...)
	digest, err := s.pullRepoTags(ctx, repo, expectedTags...)
	assert.NoError(err)
	assert.Equal(expectedDigest, digest)
}
