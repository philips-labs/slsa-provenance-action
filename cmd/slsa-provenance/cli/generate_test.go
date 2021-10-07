package cli_test

import (
	"context"
	"testing"

	"github.com/philips-labs/slsa-provenance-action/cmd/slsa-provenance/cli"
)

func TestErrors(t *testing.T) {
	testCases := []struct {
		name          string
		errorMessage  string
		validArgument []string
	}{
		{
			name:         "test artifact path argument error",
			errorMessage: "no value found for required flag: -artifact_path",
			// TODO: Refactor from validArguments to "Arguments" - Allows testing for bad/happy flow
			validArgument: []string{"-artifact_path", "artifact/path"},
		},
		{
			name:          "test output path argument error",
			errorMessage:  "no value found for required flag: -output_path",
			validArgument: []string{"-output_path", "output/path"},
		},
		{
			name:          "test github context argument error",
			errorMessage:  "no value found for required flag: -github_context",
			validArgument: []string{"-github_context", "gh-context"},
		},
		{
			name:          "test runner context argument error",
			errorMessage:  "no value found for required flag: -runner_context",
			validArgument: []string{"-runner_context", "runner-context"},
		},
		// TODO: Add test case that provides no arguments (use make slice function)
	}

	cli := cli.Generate()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := cli.Parse(tc.validArgument)
			if err != nil {
				t.Error("Failed parsing arguments")
			}
			err = cli.Run(context.Background())

			// TODO: Happy flow: If err message is empty string, go test happy flow

			if err == nil {
				t.Error("Expected an error but did not generate one")
			} else {
				if err.Error() != tc.errorMessage {
					t.Errorf("Expected error to match: %s, got: %v", tc.errorMessage, err)
				}
			}
		})
	}

}
