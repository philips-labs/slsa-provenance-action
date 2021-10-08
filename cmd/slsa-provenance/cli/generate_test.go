package cli_test

import (
	"context"
	"fmt"
	"testing"

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

	cli := cli.Generate()
	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			err := cli.ParseAndRun(context.Background(), tc.arguments)

			if tc.err != nil {
				if err == nil {
					tt.Error("Expected an error but did not generate one")
				} else {
					if err.Error() != tc.err.Error() {
						tt.Errorf("Expected error to match: %v, got: %v", tc.err, err)
					}
				}
			} else {
				fmt.Println("Add happyflow tests")
			}
		})
	}

}
