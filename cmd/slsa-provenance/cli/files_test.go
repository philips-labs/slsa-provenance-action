package cli_test

import (
	"fmt"
	"os"
	"path"
	"runtime"
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
			err:       cli.RequiredFlagError("artifact-path"),
			arguments: make([]string, 0),
		},
		{
			name: "only providing --artifact-path",
			err:  cli.RequiredFlagError("github-context"),
			arguments: []string{
				"--artifact-path",
				path.Join(rootDir, "bin/slsa-provenance"),
			},
		},
		{
			name: "without -runner-context",
			err:  cli.RequiredFlagError("runner-context"),
			arguments: []string{
				"--artifact-path",
				path.Join(rootDir, "bin/slsa-provenance"),
				"--github-context",
				githubContext,
				"--output-path",
				provenanceFile,
			},
		},
		{
			name: "invalid --artifact-path",
			err:  fmt.Errorf("failed to generate provenance: resource path not found: [provided=non-existing-folder/unknown-file]"),
			arguments: []string{
				"--artifact-path",
				"non-existing-folder/unknown-file",
				"--github-context",
				githubContext,
				"--runner-context",
				runnerContext,
			},
		},
		{
			name: "all arguments explicit",
			err:  nil,
			arguments: []string{
				"--artifact-path",
				path.Join(rootDir, "bin/slsa-provenance"),
				"--github-context",
				githubContext,
				"--output-path",
				provenanceFile,
				"--runner-context",
				runnerContext,
			},
		},
		{
			name: "With extra materials",
			err:  nil,
			arguments: []string{
				"--artifact-path",
				path.Join(rootDir, "bin/slsa-provenance"),
				"--github-context",
				githubContext,
				"--output-path",
				provenanceFile,
				"--runner-context",
				runnerContext,
				"--extra-materials",
				path.Join(rootDir, "test-data/materials-valid.json"),
			},
		},
		{
			name: "With broken extra materials",
			err:  fmt.Errorf("failed retrieving extra materials for %s: unexpected EOF", path.Join(rootDir, "test-data/materials-broken.not-json")),
			arguments: []string{
				"--artifact-path",
				path.Join(rootDir, "bin/slsa-provenance"),
				"--github-context",
				githubContext,
				"--output-path",
				provenanceFile,
				"--runner-context",
				runnerContext,
				"--extra-materials",
				path.Join(rootDir, "test-data/materials-broken.not-json"),
			},
		},
		{
			name: "With non-existent extra materials",
			err:  fmt.Errorf("failed retrieving extra materials: open %s: no such file or directory", "non-existing-folder/unknown-file"),
			arguments: []string{
				"--artifact-path",
				path.Join(rootDir, "bin/slsa-provenance"),
				"--github-context",
				githubContext,
				"--output-path",
				provenanceFile,
				"--runner-context",
				runnerContext,
				"--extra-materials",
				fmt.Sprintf("%s,%s", path.Join(rootDir, "test-data/materials-valid.json"), "non-existing-folder/unknown-file"),
			},
		},
		{
			name: "With broken extra materials (no uri)",
			err:  fmt.Errorf("failed retrieving extra materials for %s: empty or missing \"uri\" for material", path.Join(rootDir, "test-data/materials-no-uri.json")),
			arguments: []string{
				"--artifact-path",
				path.Join(rootDir, "bin/slsa-provenance"),
				"--github-context",
				githubContext,
				"--output-path",
				provenanceFile,
				"--runner-context",
				runnerContext,
				"--extra-materials",
				path.Join(rootDir, "test-data/materials-no-uri.json"),
			},
		},
		{
			name: "With broken extra materials (no digest)",
			err:  fmt.Errorf("failed retrieving extra materials for %s: empty or missing \"digest\" for material", path.Join(rootDir, "test-data/materials-no-digest.json")),
			arguments: []string{
				"--artifact-path",
				path.Join(rootDir, "bin/slsa-provenance"),
				"--github-context",
				githubContext,
				"--output-path",
				provenanceFile,
				"--runner-context",
				runnerContext,
				"--extra-materials",
				path.Join(rootDir, "test-data/materials-no-digest.json"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			assert := assert.New(tt)

			output, err := executeCommand(cli.Files(), tc.arguments...)
			defer func() {
				_ = os.Remove(provenanceFile)
			}()

			if tc.err != nil {
				assert.EqualError(err, tc.err.Error())
			} else {
				assert.NoError(err)
				assert.Contains(output, "Saving provenance to")
				if assert.FileExists(provenanceFile) {
					content, err := os.ReadFile(provenanceFile)
					assert.NoError(err)
					assert.Greater(len(content), 1)
				}
			}
		})
	}
}
