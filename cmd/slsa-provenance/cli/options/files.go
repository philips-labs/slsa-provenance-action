package options

import (
	"fmt"

	"github.com/spf13/cobra"
)

// RequiredFlagError creates a required flag error for the given flag name
func RequiredFlagError(flagName string) error {
	return fmt.Errorf("no value found for required flag: %s", flagName)
}

// FilesOptions Commandline flags used for the generate files command.
type FilesOptions struct {
	GenerateOptions
	ArtifactPath string
}

// GetArtifactPath The file(s) or directory of artifacts to include in provenance.
func (o *FilesOptions) GetArtifactPath() (string, error) {
	if o.ArtifactPath == "" {
		return "", RequiredFlagError("artifact-path")
	}
	return o.ArtifactPath, nil
}

// AddFlags Registers the flags with the cobra.Command.
func (o *FilesOptions) AddFlags(cmd *cobra.Command) {
	o.GenerateOptions.AddFlags(cmd)
	cmd.PersistentFlags().StringVar(&o.ArtifactPath, "artifact-path", "", "The file(s) or directory of artifacts to include in provenance.")
}
