package cli_test

import (
	"context"
	"testing"

	"github.com/philips-labs/slsa-provenance-action/cmd/slsa-provenance/cli"
)

const (
	ArtifactPathError = "no value found for required flag: -artifact_path"
)

func TestErrors(t *testing.T) {
	cli := cli.Generate()
	err := cli.Exec(context.Background(), make([]string, 0))

	if err == nil {
		t.Error("Expected an error but did not generate one")
	} else {
		if err.Error() != ArtifactPathError {
			t.Errorf("Expected error to match: %s", ArtifactPathError)
		}
	}
}
