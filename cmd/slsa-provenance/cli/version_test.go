package cli_test

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/philips-labs/slsa-provenance-action/cmd/slsa-provenance/cli"
)

func TestVersionCliText(t *testing.T) {
	assert := assert.New(t)

	expected := fmt.Sprintf(`GitVersion:    devel
GitCommit:     unknown
GitTreeState:  unknown
BuildDate:     unknown
GoVersion:     %s
Compiler:      %s
Platform:      %s/%s

`, runtime.Version(), runtime.Compiler, runtime.GOOS, runtime.GOARCH)

	output, err := executeCommand(cli.Version())
	assert.NoError(err)
	assert.Equal(expected, output)
}

func TestVersionCliJSON(t *testing.T) {
	assert := assert.New(t)

	expected := fmt.Sprintf(`{
  "git_version": "devel",
  "git_commit": "unknown",
  "git_tree_state": "unknown",
  "build_date": "unknown",
  "go_version": "%s",
  "compiler": "%s",
  "platform": "%s/%s"
}
`, runtime.Version(), runtime.Compiler, runtime.GOOS, runtime.GOARCH)

	output, err := executeCommand(cli.Version(), "--json")
	assert.NoError(err)
	assert.Equal(expected, output)
}
