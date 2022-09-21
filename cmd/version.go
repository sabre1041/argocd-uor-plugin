package cmd

import (
	"fmt"
	"html/template"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	// commit is the head commit from git
	commit string
	// buildDate in ISO8601 format
	buildDate string
	// version describes the version of the client
	// set at build time or detected during runtime.
	version = "v0.0.0-unknown"
	// buildData set at build time to add extra information
	// to the version.
	buildData string
)

var versionTemplate = `argocd-uor-plugin:

 Version:	{{ .Version }}
 Go Version:	{{ .GoVersion }}
 Git Commit:	{{ .GitCommit }}
 Build Date:	{{ .BuildDate }}
 Platform:	{{ .Platform }}
`

type versionInfo struct {
	Platform  string
	Version   string
	GitCommit string
	GoVersion string
	BuildDate string
}

// NewVersionCommand initializes the version command
func NewVersionCommand() *cobra.Command {

	// versionCmd represents the version command
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print argocd-uor-plugin version",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return getVersion(cmd)
		},
	}

	return versionCmd
}

// getVersion will output the templated version message.
func getVersion(cmd *cobra.Command) error {

	versionWithBuild := func() string {
		if buildData != "" {
			return fmt.Sprintf("%s+%s", version, buildData)
		}

		return version
	}

	versionInfo := versionInfo{
		Version:   versionWithBuild(),
		GitCommit: commit,
		BuildDate: buildDate,
		GoVersion: runtime.Version(),
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}

	tmp, err := template.New("version").Parse(versionTemplate)
	if err != nil {
		return fmt.Errorf("template parsing error: %v", err)
	}

	return tmp.Execute(cmd.OutOrStdout(), versionInfo)
}
