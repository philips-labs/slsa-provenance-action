package oci

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/docker/docker/api/types"
)

func (c *ContainerSubjecter) pullRepoTags(ctx context.Context, repo string, tags ...string) (string, error) {
	if len(tags) == 0 {
		tags = append(tags, "latest")
	}

	var digest string
	for _, t := range tags {
		img := fmt.Sprintf("%s:%s", repo, t)
		reader, err := c.cli.ImagePull(ctx, img, types.ImagePullOptions{})
		if err != nil {
			return "", err
		}
		defer reader.Close()
		digest, err = grepDigest(reader)
		if err != nil {
			return "", err
		}
	}

	return digest, nil
}

func grepDigest(reader io.Reader) (string, error) {
	scanner := bufio.NewScanner(reader)
	r, err := regexp.Compile(`Digest: (.*)"}`)
	if err != nil {
		return "", err
	}

	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "Digest:") {
			m := r.FindStringSubmatch(line)
			if len(m) > 1 {
				return m[1], nil
			}
		}
	}

	return "", fmt.Errorf("digest not found")
}
