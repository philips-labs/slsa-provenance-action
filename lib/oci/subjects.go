package oci

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/docker/docker/client"

	"github.com/philips-labs/slsa-provenance-action/lib/intoto"
)

// ContainerSubjecter implements Subjector to retrieve Subject from given container
// if digest is given, it will also compare matches with the given digest
type ContainerSubjecter struct {
	cli    *client.Client
	repo   string
	digest string
	tags   []string
}

// NewContainerSubjecter walks the docker tags to retrieve the digests.
// If digest is non empty string it will be used to compare the rerieved digest
// to match the given digest
func NewContainerSubjecter(cli *client.Client, repo, digest string, tags ...string) *ContainerSubjecter {
	return &ContainerSubjecter{cli, repo, digest, tags}
}

// Subjects walks the file or directory at "root" and hashes all files.
func (c *ContainerSubjecter) Subjects() ([]intoto.Subject, error) {
	digest, err := c.pullRepoTags(context.TODO(), c.repo, c.tags...)
	if err != nil {
		return nil, err
	}
	if c.digest != "" && c.digest != digest {
		return nil, fmt.Errorf("did not get expected digest, got %s, expected %s", digest, c.digest)
	}
	digestParts := strings.Split(digest, ":")
	subjects := make([]intoto.Subject, len(c.tags))

	for i, t := range c.tags {
		subjects[i] = intoto.Subject{
			Name:   fmt.Sprintf("%s:%s", c.repo, t),
			Digest: intoto.DigestSet{digestParts[0]: digestParts[1]},
		}
	}

	return subjects, nil
}

// ShaSum256HexEncoded calculates a SHA256 checksum from the content
func ShaSum256HexEncoded(b []byte) string {
	sha := sha256.Sum256(b)
	shaHex := hex.EncodeToString(sha[:])

	return shaHex
}
