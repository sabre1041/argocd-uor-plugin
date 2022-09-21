package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	genericclioptions "k8s.io/cli-runtime/pkg/genericclioptions"
)

type PluginOptions struct {
	IOStreams genericclioptions.IOStreams
}

// NewRootCommand returns a new instance of the root command
func NewRootCommand() *cobra.Command {

	var rootCmd = &cobra.Command{
		Use:   "argocd-uor-plugin",
		Short: "Argo CD UOR Plugin",
		Long:  `Argo CD plugin to generate Kubernetes resources from UOR collections`,
	}

	pluginOptions := &PluginOptions{
		IOStreams: genericclioptions.IOStreams{
			In:     os.Stdin,
			Out:    os.Stdout,
			ErrOut: os.Stderr,
		},
	}

	rootCmd.AddCommand(NewGenerateCommand(pluginOptions))
	rootCmd.AddCommand(NewVersionCommand())

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	return rootCmd
}

func init() {
	viper.SetEnvPrefix("ARGOCD_ENV")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	viper.AutomaticEnv()

}
