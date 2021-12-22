package cli_test

import (
	"encoding/base64"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/philips-labs/slsa-provenance-action/cmd/slsa-provenance/cli"
)

func TestGenerateContainerCliOptions(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	provenanceFile := path.Join(path.Dir(filename), "provenance.json")

	base64GitHubContext := base64.StdEncoding.EncodeToString([]byte(githubContext))
	base64RunnerContext := base64.StdEncoding.EncodeToString([]byte(runnerContext))

	testCases := []struct {
		name      string
		err       error
		arguments []string
	}{
		{
			name:      "without commandline flags",
			err:       cli.RequiredFlagError("github-context"),
			arguments: make([]string, 0),
		},
		{
			name: "only github-context given",
			err:  cli.RequiredFlagError("runner-context"),
			arguments: []string{
				"--github-context",
				base64GitHubContext,
			},
		},
		{
			name: "only context flags given",
			err:  cli.RequiredFlagError("repository"),
			arguments: []string{
				"--github-context",
				base64GitHubContext,
				"--runner-context",
				base64RunnerContext,
			},
		},
		{
			name: "contexts and tags given",
			err:  cli.RequiredFlagError("repository"),
			arguments: []string{
				"--github-context",
				base64GitHubContext,
				"--runner-context",
				base64RunnerContext,
				"--tags",
				"v0.4.0,33ba3da2213c83ce02df0f2f6ba925ec79037f9d",
			},
		},
		{
			name: "contexts, repo and tags given",
			err:  cli.RequiredFlagError("digest"),
			arguments: []string{
				"--github-context",
				base64GitHubContext,
				"--runner-context",
				base64RunnerContext,
				"--repository",
				"ghcr.io/philips-labs/slsa-provenance",
				"--tags",
				"v0.4.0,33ba3da2213c83ce02df0f2f6ba925ec79037f9d",
			},
		},
		{
			name: "all flags given",
			err:  nil,
			arguments: []string{
				"--github-context",
				base64GitHubContext,
				"--runner-context",
				base64RunnerContext,
				"--repository",
				"ghcr.io/philips-labs/slsa-provenance",
				"--tags",
				"v0.4.0,33ba3da2213c83ce02df0f2f6ba925ec79037f9d",
				"--digest",
				"sha256:194b471a878add368bf02a7935fa099024576c029491bcefaeb87f81efa093a3",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			assert := assert.New(tt)

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
