package options

import (
	"github.com/spf13/cobra"
)

type GitHubReleaseOptions struct {
	GenerateOptions
	ArtifactPath string
	OutputPath   string
	TagName      string
}

func (o *GitHubReleaseOptions) GetArtifactPath() (string, error) {
	if o.ArtifactPath == "" {
		return "", RequiredFlagError("artifact-path")
	}
	return o.ArtifactPath, nil
}

func (o *GitHubReleaseOptions) GetOutputPath() (string, error) {
	if o.ArtifactPath == "" {
		return "", RequiredFlagError("output-path")
	}
	return o.OutputPath, nil
}

func (o *GitHubReleaseOptions) GetTagName() (string, error) {
	if o.TagName == "" {
		return "", RequiredFlagError("tag-name")
	}
	return o.TagName, nil
}

func (o *GitHubReleaseOptions) AddFlags(cmd *cobra.Command) {
	o.GenerateOptions.AddFlags(cmd)
	cmd.PersistentFlags().StringVar(&o.ArtifactPath, "artifact-path", "", "The file(s) or directory of artifacts to include in provenance.")
	cmd.PersistentFlags().StringVar(&o.OutputPath, "output-path", "provenance.json", "The path to which the generated provenance should be written.")
	cmd.PersistentFlags().StringVar(&o.TagName, "tag-name", "", `The github release to generate provenance on.
	(if set the artifacts will be downloaded from the release and the provenance will be added as an additional release asset.)`)
}
