package oci

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/philips-labs/slsa-provenance-action/pkg/intoto"
)

func TestSubjects(t *testing.T) {
	assert := assert.New(t)

	repo := "ghcr.io/philips-labs/slsa-provenance"

	opts := WithDefaultClientOptions(context.Background(), false, false)

	errorCases := []struct {
		name   string
		repo   string
		tags   []string
		digest string
		err    string
	}{
		{
			name:   "without arguments",
			repo:   "",
			tags:   nil,
			digest: "",
			err:    "parsing reference \":latest\": could not parse reference: :latest",
		},
		{
			name:   "with non existing tag",
			repo:   repo,
			tags:   []string{"non-existing"},
			digest: "",
			err:    "GET https://ghcr.io/v2/philips-labs/slsa-provenance/manifests/non-existing: MANIFEST_UNKNOWN: manifest unknown",
		},
		{
			name:   "invalid digest",
			repo:   repo,
			tags:   []string{"v0.4.0"},
			digest: "sha256:284b471a878add368bf02a7935fa099024576c029491bcefaeb87f81efa093a4",
			err:    "did not get expected digest, got sha256:194b471a878add368bf02a7935fa099024576c029491bcefaeb87f81efa093a3, expected sha256:284b471a878add368bf02a7935fa099024576c029491bcefaeb87f81efa093a4",
		},
	}

	happyCases := []struct {
		name   string
		tags   []string
		digest string
		count  int
	}{
		{
			name:   "single tag (git tag)",
			tags:   []string{"v0.4.0"},
			digest: "sha256:194b471a878add368bf02a7935fa099024576c029491bcefaeb87f81efa093a3",
			count:  1,
		}, {
			name:   "single tag (commit hash)",
			tags:   []string{"33ba3da2213c83ce02df0f2f6ba925ec79037f9d"},
			digest: "sha256:194b471a878add368bf02a7935fa099024576c029491bcefaeb87f81efa093a3",
			count:  1,
		}, {
			name:   "muliple tags",
			tags:   []string{"v0.4.0", "33ba3da2213c83ce02df0f2f6ba925ec79037f9d"},
			digest: "sha256:194b471a878add368bf02a7935fa099024576c029491bcefaeb87f81efa093a3",
			count:  2,
		},
	}

	for _, tc := range happyCases {
		t.Run(tc.name, func(tt *testing.T) {
			subjecter := NewContainerSubjecter(repo, tc.digest, tc.tags, opts...)
			s, err := subjecter.Subjects()
			assert.NoError(err)
			assert.NotNil(s)
			assert.Len(s, tc.count)

			for i := 0; i < tc.count; i++ {
				assertSubject(assert, s, repo, tc.tags[i], tc.digest)
			}
		})
	}

	for _, tc := range errorCases {
		t.Run(tc.name, func(tt *testing.T) {
			subjecter := NewContainerSubjecter(tc.repo, tc.digest, tc.tags, opts...)
			s, err := subjecter.Subjects()
			assert.EqualError(err, tc.err)
			assert.Nil(s)
		})
	}
}

func assertSubject(assert *assert.Assertions, subject []intoto.Subject, repo, tag, digest string) {
	subjectName := fmt.Sprintf("%s:%s", repo, tag)
	digestValue := strings.Split(digest, ":")[1]
	assert.Contains(subject, intoto.Subject{Name: subjectName, Digest: intoto.DigestSet{"sha256": digestValue}})
}
