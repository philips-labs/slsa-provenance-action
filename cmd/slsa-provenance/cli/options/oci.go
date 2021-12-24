package options

import (
	"github.com/spf13/cobra"
)

// OCIOptions Commandline flags used for the generate oci command.
type OCIOptions struct {
	GenerateOptions
	Repository string
	Digest     string
	Tags       []string
}

// GetRepository The oci repository to search for the given tags.
func (o *OCIOptions) GetRepository() (string, error) {
	if o.Repository == "" {
		return "", RequiredFlagError("repository")
	}
	return o.Repository, nil
}

// GetDigest The digest to validate the tag digests against.
func (o *OCIOptions) GetDigest() (string, error) {
	if o.Digest == "" {
		return "", RequiredFlagError("digest")
	}
	return o.Digest, nil
}

// GetTags The tags to add as provenance subjects.
func (o *OCIOptions) GetTags() ([]string, error) {
	return o.Tags, nil
}

// AddFlags Registers the flags with the cobra.Command.
func (o *OCIOptions) AddFlags(cmd *cobra.Command) {
	o.GenerateOptions.AddFlags(cmd)
	cmd.PersistentFlags().StringVar(&o.Repository, "repository", "", "The repository of the oci artifact.")
	cmd.PersistentFlags().StringVar(&o.Digest, "digest", "", "The digest for the oci artifact.")
	cmd.PersistentFlags().StringSliceVar(&o.Tags, "tags", []string{"latest"}, "The given tags for this oci release.")
}
