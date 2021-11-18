package cli_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/philips-labs/slsa-provenance-action/cmd/slsa-provenance/cli"
)

func TestGenerate(t *testing.T) {
	assert := assert.New(t)

	cmd := cli.Generate()

	assert.Len(cmd.Commands(), 3)
	output, err := executeCommand(cmd)

	assert.NoError(err)
	assert.Contains(output, "Generate provenance using subcommands\n\nUsage:\n")
}
