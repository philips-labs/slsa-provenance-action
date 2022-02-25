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

const (
	unknownFile = "non-existing-folder/unknown-file"
)

func TestGenerateFilesCliOptions(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	rootDir := path.Join(path.Dir(filename), "../../..")
	provenanceFile := path.Join(rootDir, "bin/unittest-provenance.json")

	testCases := []struct {
		name        string
		err         error
		arguments   []string
		environment map[string]string
	}{
		{
			name:        "no environment variables",
			err:         cli.RequiredEnvironmentVariableError("GITHUB_CONTEXT"),
			arguments:   []string{},
			environment: map[string]string{},
		},
		{
			name:      "only github-context given",
			err:       cli.RequiredEnvironmentVariableError("RUNNER_CONTEXT"),
			arguments: []string{},
			environment: map[string]string{
				"GITHUB_CONTEXT": githubContext,
			},
		},
		{
			name:      "only contexts given",
			err:       cli.RequiredFlagError("artifact-path"),
			arguments: []string{},
			environment: map[string]string{
				"GITHUB_CONTEXT": githubContext,
				"RUNNER_CONTEXT": runnerContext,
			},
		},
		{
			name: "invalid --artifact-path",
			err:  fmt.Errorf("failed to generate provenance: lstat non-existing-folder/unknown-file: no such file or directory"),
			arguments: []string{
				"--artifact-path",
				unknownFile,
			},
			environment: map[string]string{
				"GITHUB_CONTEXT": githubContext,
				"RUNNER_CONTEXT": runnerContext,
			},
		},
		{
			name: "invalid --output-path",
			err:  fmt.Errorf("no value found for required flag: output-path"),
			arguments: []string{
				"--artifact-path",
				unknownFile,
				"--output-path",
				"",
			},
			environment: map[string]string{
				"GITHUB_CONTEXT": githubContext,
				"RUNNER_CONTEXT": runnerContext,
			},
		},
		{
			name: "all arguments explicit",
			err:  nil,
			arguments: []string{
				"--artifact-path",
				path.Join(rootDir, "bin/slsa-provenance"),
				"--output-path",
				provenanceFile,
			},
			environment: map[string]string{
				"GITHUB_CONTEXT": githubContext,
				"RUNNER_CONTEXT": runnerContext,
			},
		},
		{
			name: "With extra materials",
			err:  nil,
			arguments: []string{
				"--artifact-path",
				path.Join(rootDir, "bin/slsa-provenance"),
				"--output-path",
				provenanceFile,
				"--extra-materials",
				path.Join(rootDir, "test-data/materials-valid.json"),
			},
			environment: map[string]string{
				"GITHUB_CONTEXT": githubContext,
				"RUNNER_CONTEXT": runnerContext,
			},
		},
		{
			name: "With broken extra materials",
			err:  fmt.Errorf("failed retrieving extra materials for %s: unexpected EOF", path.Join(rootDir, "test-data/materials-broken.not-json")),
			arguments: []string{
				"--artifact-path",
				path.Join(rootDir, "bin/slsa-provenance"),
				"--output-path",
				provenanceFile,
				"--extra-materials",
				path.Join(rootDir, "test-data/materials-broken.not-json"),
			},
			environment: map[string]string{
				"GITHUB_CONTEXT": githubContext,
				"RUNNER_CONTEXT": runnerContext,
			},
		},
		{
			name: "With non-existent extra materials",
			err:  fmt.Errorf("failed retrieving extra materials: open %s: no such file or directory", unknownFile),
			arguments: []string{
				"--artifact-path",
				path.Join(rootDir, "bin/slsa-provenance"),
				"--output-path",
				provenanceFile,
				"--extra-materials",
				fmt.Sprintf("%s,%s", path.Join(rootDir, "test-data/materials-valid.json"), unknownFile),
			},
			environment: map[string]string{
				"GITHUB_CONTEXT": githubContext,
				"RUNNER_CONTEXT": runnerContext,
			},
		},
		{
			name: "With broken extra materials (no uri)",
			err:  fmt.Errorf("failed retrieving extra materials for %s: empty or missing \"uri\" for material", path.Join(rootDir, "test-data/materials-no-uri.json")),
			arguments: []string{
				"--artifact-path",
				path.Join(rootDir, "bin/slsa-provenance"),
				"--output-path",
				provenanceFile,
				"--extra-materials",
				path.Join(rootDir, "test-data/materials-no-uri.json"),
			},
			environment: map[string]string{
				"GITHUB_CONTEXT": githubContext,
				"RUNNER_CONTEXT": runnerContext,
			},
		},
		{
			name: "With broken extra materials (no digest)",
			err:  fmt.Errorf("failed retrieving extra materials for %s: empty or missing \"digest\" for material", path.Join(rootDir, "test-data/materials-no-digest.json")),
			arguments: []string{
				"--artifact-path",
				path.Join(rootDir, "bin/slsa-provenance"),
				"--output-path",
				provenanceFile,
				"--extra-materials",
				path.Join(rootDir, "test-data/materials-no-digest.json"),
			},
			environment: map[string]string{
				"GITHUB_CONTEXT": githubContext,
				"RUNNER_CONTEXT": runnerContext,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			assert := assert.New(tt)

			// Set environment
			os.Clearenv()
			for k, v := range tc.environment {
				os.Setenv(k, v)
			}

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
