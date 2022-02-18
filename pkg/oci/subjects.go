package oci

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/google/go-containerregistry/pkg/crane"

	"github.com/philips-labs/slsa-provenance-action/pkg/intoto"
)

// ContainerSubjecter implements Subjector to retrieve Subject from given container
// if digest is given, it will also compare matches with the given digest
type ContainerSubjecter struct {
	options []crane.Option
	repo    string
	digest  string
	tags    []string
}

// NewContainerSubjecter walks the docker tags to retrieve the digests.
// If digest is non empty string it will be used to compare the rerieved digest
// to match the given digest
func NewContainerSubjecter(repo, digest string, tags []string, options ...crane.Option) *ContainerSubjecter {
	return &ContainerSubjecter{options, repo, digest, tags}
}

// Subjects walks the file or directory at "root" and hashes all files.
func (c *ContainerSubjecter) Subjects() ([]intoto.Subject, error) {
	subjects := make([]intoto.Subject, len(c.tags))

	if c.tags == nil || len(c.tags) == 0 {
		c.tags = []string{"latest"}
	}

	for i, t := range c.tags {
		digest, err := crane.Digest(fmt.Sprintf("%s:%s", c.repo, t), c.options...)
		if err != nil {
			return nil, err
		}
		if c.digest != "" && c.digest != digest {
			return nil, fmt.Errorf("did not get expected digest, got %s, expected %s", digest, c.digest)
		}
		digestParts := strings.Split(digest, ":")
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
