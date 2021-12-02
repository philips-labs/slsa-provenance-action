package cli

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// Version creates an instance of *cobra.Command to print version info
func Version() *cobra.Command {
	var outputJSON bool

	cmd := &cobra.Command{
		Use:   "version",
		Short: fmt.Sprintf("Prints the %s version", cliName),
		Long:  fmt.Sprintf("Prints the %s version", cliName),
		RunE: func(cmd *cobra.Command, args []string) error {
			v := VersionInfo()
			res := v.String()
			if outputJSON {
				j, err := v.JSONString()
				if err != nil {
					return fmt.Errorf("unable to generate JSON from version info: %w", err)
				}
				res = j
			}
			fmt.Fprintln(cmd.OutOrStdout(), res)
			return nil
		},
	}

	cmd.Flags().BoolVar(&outputJSON, "json", false, "print version information as JSON")

	return cmd
}

// Base version information.
//
// This is the fallback data used when version information from git is not
// provided via go ldflags (e.g. via Makefile).
var (
	// Output of "git describe". The prerequisite is that the branch should be
	// tagged using the correct versioning strategy.
	GitVersion string = "devel"
	// SHA1 from git, output of $(git rev-parse HEAD)
	gitCommit = "unknown"
	// State of git tree, either "clean" or "dirty"
	gitTreeState = "unknown"
	// Build date in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ')
	buildDate = "unknown"
)

// Info holds the version information of the binary
type Info struct {
	GitVersion   string `json:"git_version,omitempty"`
	GitCommit    string `json:"git_commit,omitempty"`
	GitTreeState string `json:"git_tree_state,omitempty"`
	BuildDate    string `json:"build_date,omitempty"`
	GoVersion    string `json:"go_version,omitempty"`
	Compiler     string `json:"compiler,omitempty"`
	Platform     string `json:"platform,omitempty"`
}

// VersionInfo creates an instance of the Info structure
func VersionInfo() Info {
	return Info{
		GitVersion:   GitVersion,
		GitCommit:    gitCommit,
		GitTreeState: gitTreeState,
		BuildDate:    buildDate,
		GoVersion:    runtime.Version(),
		Compiler:     runtime.Compiler,
		Platform:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

// String returns the string representation of the version info
func (i *Info) String() string {
	b := strings.Builder{}
	w := tabwriter.NewWriter(&b, 0, 0, 2, ' ', 0)

	fmt.Fprintf(w, "GitVersion:\t%s\n", i.GitVersion)
	fmt.Fprintf(w, "GitCommit:\t%s\n", i.GitCommit)
	fmt.Fprintf(w, "GitTreeState:\t%s\n", i.GitTreeState)
	fmt.Fprintf(w, "BuildDate:\t%s\n", i.BuildDate)
	fmt.Fprintf(w, "GoVersion:\t%s\n", i.GoVersion)
	fmt.Fprintf(w, "Compiler:\t%s\n", i.Compiler)
	fmt.Fprintf(w, "Platform:\t%s\n", i.Platform)

	w.Flush()
	return b.String()
}

// JSONString returns the JSON representation of the version info
func (i *Info) JSONString() (string, error) {
	b, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		return "", err
	}

	return string(b), nil
}
