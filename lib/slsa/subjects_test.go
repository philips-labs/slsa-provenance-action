package slsa

import (
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/philips-labs/slsa-provenance-action/lib/intoto"
)

func TestSubjects(t *testing.T) {
	assert := assert.New(t)

	s, err := subjects("/invalid-path")
	assert.Error(err)
	assert.Nil(s)

	_, filename, _, _ := runtime.Caller(0)
	rootDir := path.Join(path.Dir(filename), "../..")
	artifactPath := path.Join(rootDir, "bin")
	binaryName := "slsa-provenance"
	binaryPath := path.Join(artifactPath, binaryName)

	s, err = subjects(artifactPath)
	assert.NoError(err)
	assert.NotNil(s)
	assert.Len(s, 1)
	AssertSubject(assert, s, binaryName, binaryPath)

	s, err = subjects(binaryPath)
	assert.NoError(err)
	assert.NotNil(s)
	assert.Len(s, 1)
	AssertSubject(assert, s, binaryName, binaryPath)

	s, err = subjects(".")
	assert.NoError(err)
	assert.NotNil(s)

	assert.Len(s, 4)
	AssertSubject(assert, s, "provenance_test.go", path.Join(".", "provenance_test.go"))
	AssertSubject(assert, s, "provenance.go", path.Join(".", "provenance.go"))
	AssertSubject(assert, s, "subjects_test.go", path.Join(".", "subjects_test.go"))
	AssertSubject(assert, s, "subjects.go", path.Join(".", "subjects.go"))
}

func AssertSubject(assert *assert.Assertions, subject []intoto.Subject, binaryName, binaryPath string) {
	binary, err := os.ReadFile(binaryPath)
	if !assert.NoError(err) {
		return
	}

	shaHex := ShaSum256HexEncoded(binary)
	assert.Contains(subject, intoto.Subject{Name: binaryName, Digest: intoto.DigestSet{"sha256": shaHex}})
}
