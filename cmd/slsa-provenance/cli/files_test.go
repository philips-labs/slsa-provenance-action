package cli_test

import (
	"context"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/philips-labs/slsa-provenance-action/cmd/slsa-provenance/cli"
)

func TestGenerateFilesCliOptions(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	rootDir := path.Join(path.Dir(filename), "../../..")
	provenanceFile := path.Join(rootDir, "bin/unittest-provenance.json")

	testCases := []struct {
		name      string
		err       error
		arguments []string
	}{
		{
			name:      "without commandline flags",
			err:       cli.RequiredFlagError("-artifact_path"),
			arguments: make([]string, 0),
		},
		{
			name: "only providing artifact_path",
			err:  cli.RequiredFlagError("-github_context"),
			arguments: []string{
				"-artifact_path",
				path.Join(rootDir, "bin/slsa-provenance"),
			},
		},
		{
			name: "without runner_context",
			err:  cli.RequiredFlagError("-runner_context"),
			arguments: []string{
				"-artifact_path",
				path.Join(rootDir, "bin/slsa-provenance"),
				"-github_context",
				githubContext,
				"-output_path",
				provenanceFile,
			},
		},
		{
			name: "invalid artifact_path",
			err:  fmt.Errorf("failed to generate provenance: resource path not found: [provided=non-existing-folder/unknown-file]"),
			arguments: []string{
				"-artifact_path",
				"non-existing-folder/unknown-file",
				"-github_context",
				githubContext,
				"-runner_context",
				runnerContext,
			},
		},
		{
			name: "all arguments explicit",
			err:  nil,
			arguments: []string{
				"-artifact_path",
				path.Join(rootDir, "bin/slsa-provenance"),
				"-github_context",
				githubContext,
				"-output_path",
				provenanceFile,
				"-runner_context",
				runnerContext,
			},
		},
		{
			name: "With extra materials",
			err:  nil,
			arguments: []string{
				"-artifact_path",
				path.Join(rootDir, "bin/slsa-provenance"),
				"-github_context",
				githubContext,
				"-output_path",
				provenanceFile,
				"-runner_context",
				runnerContext,
				"-extra_materials",
				path.Join(rootDir, "test-data/materials-valid.json"),
			},
		},
		{
			name: "With broken extra materials",
			err:  fmt.Errorf("invalid JSON in extra materials file %s: unexpected end of JSON input", path.Join(rootDir, "test-data/materials-broken.not-json")),
			arguments: []string{
				"-artifact_path",
				path.Join(rootDir, "bin/slsa-provenance"),
				"-github_context",
				githubContext,
				"-output_path",
				provenanceFile,
				"-runner_context",
				runnerContext,
				"-extra_materials",
				path.Join(rootDir, "test-data/materials-broken.not-json"),
			},
		},
		{
			name: "With non-existent extra materials",
			err:  fmt.Errorf("could not load extra materials from non-existing-folder/unknown-file: open non-existing-folder/unknown-file: no such file or directory"),
			arguments: []string{
				"-artifact_path",
				path.Join(rootDir, "bin/slsa-provenance"),
				"-github_context",
				githubContext,
				"-output_path",
				provenanceFile,
				"-runner_context",
				runnerContext,
				"-extra_materials",
				fmt.Sprintf("%s %s", path.Join(rootDir, "test-data/materials-valid.json"), "non-existing-folder/unknown-file"),
			},
		},
		{
			name: "With broken extra materials (no uri)",
			err:  fmt.Errorf("empty or missing \"uri\" field in %s", path.Join(rootDir, "test-data/materials-no-uri.json")),
			arguments: []string{
				"-artifact_path",
				path.Join(rootDir, "bin/slsa-provenance"),
				"-github_context",
				githubContext,
				"-output_path",
				provenanceFile,
				"-runner_context",
				runnerContext,
				"-extra_materials",
				path.Join(rootDir, "test-data/materials-no-uri.json"),
			},
		},
		{
			name: "With broken extra materials (no digest)",
			err:  fmt.Errorf("empty or missing \"digest\" in %s", path.Join(rootDir, "test-data/materials-no-digest.json")),
			arguments: []string{
				"-artifact_path",
				path.Join(rootDir, "bin/slsa-provenance"),
				"-github_context",
				githubContext,
				"-output_path",
				provenanceFile,
				"-runner_context",
				runnerContext,
				"-extra_materials",
				path.Join(rootDir, "test-data/materials-no-digest.json"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			assert := assert.New(tt)

			sb := strings.Builder{}

			cli := cli.Files(&sb)
			err := cli.ParseAndRun(context.Background(), tc.arguments)
			defer func() {
				_ = os.Remove(provenanceFile)
			}()

			if tc.err != nil {
				assert.EqualError(err, tc.err.Error())
			} else {
				assert.NoError(err)
				assert.Contains(sb.String(), "Saving provenance to")
				if assert.FileExists(provenanceFile) {
					content, err := os.ReadFile(provenanceFile)
					assert.NoError(err)
					assert.Greater(len(content), 1)
				}
			}
		})
	}
}
