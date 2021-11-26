package options

import (
	"fmt"

	"github.com/spf13/cobra"
)

// RequiredFlagError creates a required flag error for the given flag name
func RequiredFlagError(flagName string) error {
	return fmt.Errorf("no value found for required flag: %s", flagName)
}

type FilesOptions struct {
	GenerateOptions
	ArtifactPath string
	OutputPath   string
}

func (o *FilesOptions) GetArtifactPath() (string, error) {
	if o.ArtifactPath == "" {
		return "", RequiredFlagError("artifact-path")
	}
	return o.ArtifactPath, nil
}

func (o *FilesOptions) GetOutputPath() (string, error) {
	if o.ArtifactPath == "" {
		return "", RequiredFlagError("output-path")
	}
	return o.OutputPath, nil
}

func (o *FilesOptions) AddFlags(cmd *cobra.Command) {
	o.GenerateOptions.AddFlags(cmd)
	cmd.PersistentFlags().StringVar(&o.ArtifactPath, "artifact-path", "", "The file(s) or directory of artifacts to include in provenance.")
	cmd.PersistentFlags().StringVar(&o.OutputPath, "output-path", "provenance.json", "The path to which the generated provenance should be written.")
}
