package cli_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/philips-labs/slsa-provenance-action/cmd/slsa-provenance/cli"
)

func TestCli(t *testing.T) {
	assert := assert.New(t)

	cli := cli.New()
	assert.Len(cli.Commands(), 2)
}
