package intoto

import (
	"crypto/sha256"
	"encoding/hex"
	"io/fs"
	"os"
	"path/filepath"
)

// Subjecter retrieves subjects
type Subjecter interface {
	Subjects() ([]Subject, error)
}

// FilePathSubjecter implements Subjector to retrieve Subject from filepath
type FilePathSubjecter struct {
	root string
}

// NewFilePathSubjecter walks the file or directory at "root" and hashes all files.
func NewFilePathSubjecter(root string) *FilePathSubjecter {
	return &FilePathSubjecter{root}
}

// Subjects walks the file or directory at "root" and hashes all files.
func (f *FilePathSubjecter) Subjects() ([]Subject, error) {
	var s []Subject
	return s, filepath.Walk(f.root, func(abspath string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		relpath, err := filepath.Rel(f.root, abspath)
		if err != nil {
			return err
		}
		// Note: filepath.Rel() returns "." when "root" and "abspath" point to the same file.
		if relpath == "." {
			relpath = filepath.Base(f.root)
		}

		binary, err := os.ReadFile(abspath)
		if err != nil {
			return err
		}

		shaHex := ShaSum256HexEncoded(binary)

		s = append(s, Subject{Name: relpath, Digest: DigestSet{"sha256": shaHex}})
		return nil
	})
}

// ShaSum256HexEncoded calculates a SHA256 checksum from the content
func ShaSum256HexEncoded(b []byte) string {
	sha := sha256.Sum256(b)
	shaHex := hex.EncodeToString(sha[:])

	return shaHex
}
