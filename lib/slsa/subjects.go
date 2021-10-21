package slsa

import (
	"crypto/sha256"
	"encoding/hex"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/philips-labs/slsa-provenance-action/lib/intoto"
)

// subjects walks the file or directory at "root" and hashes all files.
func subjects(root string) ([]intoto.Subject, error) {
	var s []intoto.Subject
	return s, filepath.Walk(root, func(abspath string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		relpath, err := filepath.Rel(root, abspath)
		if err != nil {
			return err
		}
		// Note: filepath.Rel() returns "." when "root" and "abspath" point to the same file.
		if relpath == "." {
			relpath = filepath.Base(root)
		}

		binary, err := os.ReadFile(abspath)
		if err != nil {
			return err
		}

		shaHex := ShaSum256HexEncoded(binary)

		s = append(s, intoto.Subject{Name: relpath, Digest: intoto.DigestSet{"sha256": shaHex}})
		return nil
	})
}

// ShaSum256HexEncoded calculates a SHA256 checksum from the content
func ShaSum256HexEncoded(b []byte) string {
	sha := sha256.Sum256(b)
	shaHex := hex.EncodeToString(sha[:])

	return shaHex
}
