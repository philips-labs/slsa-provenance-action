package github_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/philips-labs/slsa-provenance-action/pkg/github"
)

func TestMarshalGitHubContext(t *testing.T) {
	assert := assert.New(t)

	data := `{ "token": "superSecret" }`

	var ghc github.Context
	err := json.Unmarshal([]byte(data), &ghc)
	assert.NoError(err)
	assert.Equal(github.Token("superSecret"), ghc.Token)

	j, err := json.Marshal(&ghc)
	assert.NoError(err)
	assert.NotContains(string(j), `"token":"superSecret"`)
	assert.Contains(string(j), `"token":"***"`)
}
