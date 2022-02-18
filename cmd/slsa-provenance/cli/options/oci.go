package options

import (
	"context"

	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/spf13/cobra"

	"github.com/philips-labs/slsa-provenance-action/pkg/oci"
)

// OCIOptions Commandline flags used for the generate oci command.
type OCIOptions struct {
	GenerateOptions
	Repository         string
	Digest             string
	Tags               []string
	AllowInsecure      bool
	KubernetesKeychain bool
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
	cmd.Flags().BoolVar(&o.AllowInsecure, "allow-insecure", false, "whether to allow insecure connections to registries. Don't use this for anything but testing")
	cmd.Flags().BoolVar(&o.KubernetesKeychain, "k8s-keychain", false, "whether to use the kubernetes keychain instead of the default keychain (supports workload identity).")
}

// GetRegistryClientOpts sets some sane default options for crane to authenticate
// private registries
func (o *OCIOptions) GetRegistryClientOpts(ctx context.Context) []crane.Option {
	return oci.WithDefaultClientOptions(ctx, o.KubernetesKeychain, o.AllowInsecure)
}
