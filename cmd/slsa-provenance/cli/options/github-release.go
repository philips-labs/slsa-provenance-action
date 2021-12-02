package options

import (
	"github.com/spf13/cobra"
)

// GitHubReleaseOptions Commandline flags used for the generate command.
type GitHubReleaseOptions struct {
	GenerateOptions
	ArtifactPath string
	TagName      string
}

// GetArtifactPath The location to store the GitHub Release artifact
func (o *GitHubReleaseOptions) GetArtifactPath() (string, error) {
	if o.ArtifactPath == "" {
		return "", RequiredFlagError("artifact-path")
	}
	return o.ArtifactPath, nil
}

// GetTagName The name of the GitHub tag/release
func (o *GitHubReleaseOptions) GetTagName() (string, error) {
	if o.TagName == "" {
		return "", RequiredFlagError("tag-name")
	}
	return o.TagName, nil
}

// AddFlags Registers the flags with the cobra.Command.
func (o *GitHubReleaseOptions) AddFlags(cmd *cobra.Command) {
	o.GenerateOptions.AddFlags(cmd)
	cmd.PersistentFlags().StringVar(&o.ArtifactPath, "artifact-path", "", "The file(s) or directory of artifacts to include in provenance.")
	cmd.PersistentFlags().StringVar(&o.TagName, "tag-name", "", `The github release to generate provenance on.
	(if set the artifacts will be downloaded from the release and the provenance will be added as an additional release asset.)`)
}
