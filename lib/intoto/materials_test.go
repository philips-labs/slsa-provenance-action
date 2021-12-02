package intoto

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaterials(t *testing.T) {
	assert := assert.New(t)

	validMaterials := strings.NewReader(`[
	{
		"uri": "pkg:deb/debian/stunnel4@5.50-3?arch=amd64",
		"digest": {
			"sha256": "e1731ae217fcbc64d4c00d707dcead45c828c5f762bcf8cc56d87de511e096fa"
		}
	}
]`)

	nonJSON := strings.NewReader(`[
	{
		"uri": "pkg:deb/debian/stunnel4@5.50-3?arch=amd64",
		"digest": {
			"sha256": "e1731ae217fcbc64d4c00d707dcead45c828c5f762bcf8cc56d87de511e096fa"
		}
	}`)

	withoutDigest := strings.NewReader(`[
	{
		"uri": "pkg:deb/debian/stunnel4@5.50-3?arch=amd64",
		"not-digest": {
			"sha256": "e1731ae217fcbc64d4c00d707dcead45c828c5f762bcf8cc56d87de511e096fa"
		}
	}
]`)

	withoutURI := strings.NewReader(`[
{
	"digest": {
		"sha256": "e1731ae217fcbc64d4c00d707dcead45c828c5f762bcf8cc56d87de511e096fa"
	}
}
]`)

	m, err := ReadMaterials(validMaterials)
	assert.NoError(err)
	assert.Len(m, 1)

	m, err = ReadMaterials(nonJSON)
	assert.EqualError(err, "unexpected EOF")
	assert.Nil(m)

	m, err = ReadMaterials(withoutDigest)
	assert.EqualError(err, "empty or missing \"digest\" for material")
	assert.Nil(m)

	m, err = ReadMaterials(withoutURI)
	assert.EqualError(err, "empty or missing \"uri\" for material")
	assert.Nil(m)
}
