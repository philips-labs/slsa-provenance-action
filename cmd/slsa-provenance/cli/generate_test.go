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
		validArgument string
	}{
		{
			name:          "test artifact path argument error",
			errorMessage:  "no value found for required flag: -artifact_path",
			validArgument: "-artifact_path artifact/path",
		},
		{
			name:          "test output path argument error",
			errorMessage:  "no value found for required flag: -output_path",
			validArgument: "-output_path output/path",
		},
		{
			name:          "test github context argument error",
			errorMessage:  "no value found for required flag: -github_context",
			validArgument: "-github_context gh-context",
		},
		{
			name:          "test runner context argument error",
			errorMessage:  "no value found for required flag: -runner_context",
			validArgument: "-runner_context runner-context",
		},
	}

	supplementedArguments := []string{}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			cli := cli.Generate()
			err := cli.Exec(context.Background(), supplementedArguments)
			supplementedArguments = append(supplementedArguments, tt.validArgument)

			if err == nil {
				t.Error("Expected an error but did not generate one")
			} else {
				if err.Error() != tt.errorMessage {
					t.Errorf("Expected error to match: %s", tt.errorMessage)
				}
			}
		})
	}

}
