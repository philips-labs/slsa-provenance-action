package options

import (
	"crypto/ed25519"
	"encoding/hex"
	"io/ioutil"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	defaultSignOutputPath = "provenance.signed.json"
)

// SignOptions are the Commandline flags for the 'sign' command
type SignOptions struct {
	ProvenancePath string
	KeyPath        string
	Key            string
	OutputPath     string
}

// GetKey returns the key, either from the key or key-path
func (o *SignOptions) GetKey() (ed25519.PrivateKey, error) {
	if o.KeyPath != "" && o.Key != "" {
		return nil, errors.New("Both key and key-path specified")
	}
	if o.KeyPath == "" && o.Key == "" {
		return nil, errors.New("Neither key nor key-path specified")
	}
	privKeyHex := o.Key
	if o.KeyPath != "" {
		content, err := ioutil.ReadFile(o.KeyPath)
		if err != nil {
			return nil, errors.Wrapf(err, "Error reading key file %s", o.KeyPath)
		}
		privKeyHex = string(content)
	}
	privKeyData, err := hex.DecodeString(privKeyHex)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to decode key")
	}
	if len(privKeyData) != ed25519.SeedSize {
		return nil, errors.Errorf("Decoded key has wrong size, expected %d bytes, got %d", ed25519.SeedSize, len(privKeyData))
	}
	return ed25519.NewKeyFromSeed(privKeyData), nil
}

// GetOutputPath returns the path where the output be written
func (o *SignOptions) GetOutputPath() (string, error) {
	if o.OutputPath == "" {
		return "", RequiredFlagError("--output-path")
	}
	return o.OutputPath, nil
}

// GetProvenancePath returns the path to the provenance to be signed
func (o *SignOptions) GetProvenancePath() (string, error) {
	if o.ProvenancePath == "" {
		return "", RequiredFlagError("--provenance-path")
	}
	return o.ProvenancePath, nil

}

// AddFlags registers the flags with the cmd
func (o *SignOptions) AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&o.ProvenancePath, "provenance-path", defaultGenerateOutputPath, "")
	cmd.PersistentFlags().StringVar(&o.KeyPath, "key-path", "", "")
	cmd.PersistentFlags().StringVar(&o.Key, "key", "", "")
	cmd.PersistentFlags().StringVar(&o.OutputPath, "output-path", defaultSignOutputPath, "")
}
