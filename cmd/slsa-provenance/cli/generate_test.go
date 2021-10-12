package cli_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/philips-labs/slsa-provenance-action/cmd/slsa-provenance/cli"
)

func TestErrors(t *testing.T) {
	testCases := []struct {
		name      string
		err       error
		arguments []string
	}{
		{
			name:      "test artifact path argument error",
			err:       cli.RequiredFlagError("-artifact_path"),
			arguments: make([]string, 0),
		},
		{
			name: "test github context argument error",
			err:  cli.RequiredFlagError("-github_context"),
			arguments: []string{
				"-artifact_path",
				"artifact/path",
			},
		},
		{
			name: "test output path argument error",
			err:  cli.RequiredFlagError("-output_path"),
			arguments: []string{
				"-artifact_path",
				"artifact/path",
				"-github_context",
				"gh-context",
				"-runner_context",
				"runner-context",
				"-output_path",
				"''",
			},
		},
		{
			name: "test runner context argument error",
			err:  cli.RequiredFlagError("-runner_context"),
			arguments: []string{
				"-artifact_path",
				"artifact/path",
				"-github_context",
				"gh-context",
				"-output_path",
				"output/path",
			},
		},
		{
			name: "test happy flow",
			err:  nil,
			arguments: []string{
				"-artifact_path",
				"artifact/path",
				"-github_context",
				"gh-context",
				"-output_path",
				"output/path",
				"-runner_context",
				"runner-context",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			assert := assert.New(tt)
			cli := cli.Generate()
			err := cli.ParseAndRun(context.Background(), tc.arguments)

			if tc.err != nil {
				assert.EqualError(err, tc.err.Error())
			} else {
				assert.NoError(err)
			}
		})
	}

}
