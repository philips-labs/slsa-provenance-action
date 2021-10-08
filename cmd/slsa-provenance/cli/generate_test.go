package cli_test

import (
	"context"
	"testing"

	"github.com/philips-labs/slsa-provenance-action/cmd/slsa-provenance/cli"
)

type Arguments struct {
	argument string
	value    string
}

func TestErrors(t *testing.T) {
	testCases := []struct {
		name         string
		errorMessage string
		arguments    []Arguments
	}{
		{
			name:         "test artifact path argument error",
			errorMessage: "no value found for required flag: -artifact_path",
			arguments: []Arguments{
				{argument: "", value: ""},
			},
		},
		{
			name:         "test github context argument error",
			errorMessage: "no value found for required flag: -github_context",
			arguments: []Arguments{
				{argument: "-artifact_path", value: "artifact/path"},
			},
		},
		{
			name:         "test output path argument error",
			errorMessage: "no value found for required flag: -output_path",
			arguments: []Arguments{
				{argument: "-artifact_path", value: "artifact/path"},
				{argument: "-github_context", value: "gh-context"},
			},
		},
		{
			name:         "test runner context argument error",
			errorMessage: "no value found for required flag: -runner_context",
			arguments: []Arguments{
				{argument: "-artifact_path", value: "artifact/path"},
				{argument: "-github_context", value: "gh-context"},
				{argument: "-output_path", value: "output/path"},
			},
		},
		{
			name:         "test happy flow",
			errorMessage: "",
			arguments: []Arguments{
				{argument: "-artifact_path", value: "artifact/path"},
				{argument: "-github_context", value: "gh-context"},
				{argument: "-output_path", value: "output/path"},
				{argument: "-runner_context", value: "runner-context"},
			},
		},
		{
			name:         "test without arguments",
			errorMessage: "",
			arguments:    make([]Arguments, 0),
		},
	}

	cli := cli.Generate()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			args := []string{}
			for _, v := range tc.arguments {
				args = append(args, v.argument, v.value)
			}
			err := cli.Parse(args)
			if err != nil {
				t.Error("Failed parsing arguments")
			}
			err = cli.Run(context.Background())

			if tc.errorMessage == "" {
				// TODO: Happy flow: If err message is empty string, go test happy flow
			}

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
