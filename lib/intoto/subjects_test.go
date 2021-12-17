package intoto

import (
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubjects(t *testing.T) {
	assert := assert.New(t)

	fps := NewFilePathSubjecter("/invalid-path")
	s, err := fps.Subjects()
	assert.Error(err)
	assert.Nil(s)

	_, filename, _, _ := runtime.Caller(0)
	rootDir := path.Join(path.Dir(filename), "../..")
	artifactPath := path.Join(rootDir, "bin")
	binaryName := "slsa-provenance"
	binaryPath := path.Join(artifactPath, binaryName)

	fps = NewFilePathSubjecter(artifactPath)
	s, err = fps.Subjects()
	assert.NoError(err)
	assert.NotNil(s)
	assert.Len(s, 1)
	assertSubject(assert, s, binaryName, binaryPath)

	fps = NewFilePathSubjecter(binaryPath)
	s, err = fps.Subjects()
	assert.NoError(err)
	assert.NotNil(s)
	assert.Len(s, 1)
	assertSubject(assert, s, binaryName, binaryPath)

	fps = NewFilePathSubjecter(".")
	s, err = fps.Subjects()
	assert.NoError(err)
	assert.NotNil(s)

	assert.Len(s, 6)
	assertSubject(assert, s, "intoto_test.go", path.Join(".", "intoto_test.go"))
	assertSubject(assert, s, "intoto.go", path.Join(".", "intoto.go"))
	assertSubject(assert, s, "subjects_test.go", path.Join(".", "subjects_test.go"))
	assertSubject(assert, s, "subjects.go", path.Join(".", "subjects.go"))
	assertSubject(assert, s, "materials_test.go", path.Join(".", "materials_test.go"))
	assertSubject(assert, s, "materials.go", path.Join(".", "materials.go"))
}

func assertSubject(assert *assert.Assertions, subject []Subject, binaryName, binaryPath string) {
	binary, err := os.ReadFile(binaryPath)
	if !assert.NoError(err) {
		return
	}

	shaHex := ShaSum256HexEncoded(binary)
	assert.Contains(subject, Subject{Name: binaryName, Digest: DigestSet{"sha256": shaHex}})
}
