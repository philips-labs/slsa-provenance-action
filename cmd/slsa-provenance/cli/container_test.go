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

func TestGenerateContainerCliOptions(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	rootDir := path.Join(path.Dir(filename), "../../..")
	provenanceFile := path.Join(path.Dir(filename), "provenance.json")

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
			err:       cli.RequiredFlagError("repository"),
			arguments: []string{},
			environment: map[string]string{
				"GITHUB_CONTEXT": githubContext,
				"RUNNER_CONTEXT": runnerContext,
			},
		},
		{
			name: "invalid --output-path",
			err:  fmt.Errorf("no value found for required flag: output-path"),
			arguments: []string{
				"--output-path",
				"",
			},
			environment: map[string]string{
				"GITHUB_CONTEXT": githubContext,
				"RUNNER_CONTEXT": runnerContext,
			},
		},
		{
			name: "With extra materials",
			err:  cli.RequiredFlagError("repository"),
			arguments: []string{
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
		{
			name: "contexts and tags given",
			err:  cli.RequiredFlagError("repository"),
			arguments: []string{
				"--tags",
				"v0.4.0,33ba3da2213c83ce02df0f2f6ba925ec79037f9d",
			},
			environment: map[string]string{
				"GITHUB_CONTEXT": githubContext,
				"RUNNER_CONTEXT": runnerContext,
			},
		},
		{
			name: "contexts, repo and tags given",
			err:  cli.RequiredFlagError("digest"),
			arguments: []string{
				"--repository",
				"ghcr.io/philips-labs/slsa-provenance",
				"--tags",
				"v0.4.0,33ba3da2213c83ce02df0f2f6ba925ec79037f9d",
			},
			environment: map[string]string{
				"GITHUB_CONTEXT": githubContext,
				"RUNNER_CONTEXT": runnerContext,
			},
		},
		{
			name: "all flags given",
			err:  nil,
			arguments: []string{
				"--repository",
				"ghcr.io/philips-labs/slsa-provenance",
				"--tags",
				"v0.4.0,33ba3da2213c83ce02df0f2f6ba925ec79037f9d",
				"--digest",
				"sha256:194b471a878add368bf02a7935fa099024576c029491bcefaeb87f81efa093a3",
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

			output, err := executeCommand(cli.OCI(), tc.arguments...)
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
