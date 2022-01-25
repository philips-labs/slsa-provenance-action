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
	// ghcr.io/philips-internal/cloud-healthcare-poc/examples/fhir-patient-list/backend:provenance-test
	repo := "ghcr.io/philips-internal/cloud-healthcare-poc/examples/fhir-patient-list/backend"
	expectedDigest := "sha256:4672ced51303ad56195f12f087abf47b44c564ab115dc75ac9b7c6f57a677ade"
	expectedTags := []string{"provenance-test"}
	s := NewContainerSubjecter(cli, repo, expectedDigest, expectedTags...)
	digest, err := s.pullRepoTags(ctx, repo, expectedTags...)
	assert.NoError(err)
	assert.Equal(expectedDigest, digest)
}
