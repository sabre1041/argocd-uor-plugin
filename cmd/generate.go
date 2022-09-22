package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	client "github.com/uor-framework/uor-client-go/cli"
	"github.com/uor-framework/uor-client-go/config"
)

type GeneratePluginOptions struct {
	ClientPullOptions *client.PullOptions
	AttributeQuery    string
	SourceDirectory   string
	// ProcessingDirectory string
}

// NewGenerateCommand initializes the generate command
func NewGenerateCommand(pluginOptions *PluginOptions) *cobra.Command {

	// Setting level to suppress printing out unless errors occur
	logger := log.New()
	logger.Level = log.ErrorLevel

	generatePluginOptions := &GeneratePluginOptions{
		ClientPullOptions: &client.PullOptions{
			RootOptions: &client.RootOptions{
				IOStreams: pluginOptions.IOStreams,
				Logger:    logger,
			},
		},
	}

	// generateCmd represents the generate command
	var generateCmd = &cobra.Command{
		Use:   "generate <path>",
		Short: "Generate Kubernetes resources from a UOR collection",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("<path> argument required to generate manifests")
			}
			return nil
		},

		Run: func(cmd *cobra.Command, args []string) {

			// TODO: Undo manual bind for now
			generatePluginOptions.ClientPullOptions.Source = viper.GetString("collection")
			generatePluginOptions.AttributeQuery = viper.GetString("attribute-query")
			generatePluginOptions.ClientPullOptions.PlainHTTP = viper.GetBool("plain-http")
			generatePluginOptions.ClientPullOptions.Insecure = viper.GetBool("insecure")
			generatePluginOptions.ClientPullOptions.PullAll = viper.GetBool("linked-collections-pull")

			cobra.CheckErr(generatePluginOptions.Complete(args))
			cobra.CheckErr(generatePluginOptions.Run(cmd))
		},
	}

	generateCmd.Flags().StringVarP(&generatePluginOptions.AttributeQuery, "attribute-query", "a", "attribute-query.yaml", "Path to the AttributeQuery resource ")
	viper.BindPFlag("attribute-query", generateCmd.Flags().Lookup("attribute-query"))

	generateCmd.Flags().StringVarP(&generatePluginOptions.ClientPullOptions.Source, "collection", "c", "", "Collection to retrieve")
	viper.BindPFlag("collection", generateCmd.Flags().Lookup("collection"))

	// UOR Client flags
	generateCmd.Flags().BoolVarP(&generatePluginOptions.ClientPullOptions.PlainHTTP, "plain-http", "p", false, "Use plain http when contacting to external registries")
	viper.BindPFlag("plain-http", generateCmd.Flags().Lookup("plain-http"))

	generateCmd.Flags().BoolVarP(&generatePluginOptions.ClientPullOptions.Insecure, "insecure", "i", false, "Allow connections to external registries over unverified SSL certificates")
	viper.BindPFlag("insecure", generateCmd.Flags().Lookup("insecure"))

	generateCmd.Flags().BoolVarP(&generatePluginOptions.ClientPullOptions.PullAll, "linked-collections-pull", "l", false, "Retrieve Linked Collections")
	viper.BindPFlag("linked-collections-pull", generateCmd.Flags().Lookup("linked-collections-pull"))

	return generateCmd

}

func (generatePluginOptions *GeneratePluginOptions) Complete(args []string) error {

	if len(generatePluginOptions.ClientPullOptions.Source) == 0 {
		return errors.New("Collection must be provided")
	}

	absPath, _ := filepath.Abs(args[0])
	generatePluginOptions.SourceDirectory = absPath

	attributePath := path.Join(generatePluginOptions.SourceDirectory, generatePluginOptions.AttributeQuery)

	// Attempt to load attributes query
	if _, err := os.Stat(attributePath); err == nil {
		_, err := config.ReadAttributeQuery(attributePath)

		if err != nil {
			generatePluginOptions.ClientPullOptions.Logger.Errorf("Failed to parse attributes query '%s', %v", generatePluginOptions.AttributeQuery, err)
		} else {
			// Reset AttributeQuery so a processing error does not occur
			generatePluginOptions.ClientPullOptions.AttributeQuery = attributePath
		}
	}
	return nil
}

func (generatePluginOptions *GeneratePluginOptions) Run(cmd *cobra.Command) error {

	workingDir, err := ioutil.TempDir(os.TempDir(), "argocd-uor-plugin-")

	if err != nil {
		return nil
	}

	generatePluginOptions.ClientPullOptions.Output = workingDir

	// Remove working directory once complete
	defer os.RemoveAll(workingDir)

	// Create a cache dir and a output directory for retrieved assets
	cacheDir := filepath.Join(workingDir, "cache")
	outputDir := filepath.Join(workingDir, "output")
	os.MkdirAll(cacheDir, 0750)
	os.MkdirAll(outputDir, 0750)

	generatePluginOptions.ClientPullOptions.Output = outputDir

	//TODO Set Cache Directory

	err = generatePluginOptions.ClientPullOptions.Run(cmd.Context())

	if err != nil {
		return err
	}

	// List Applicable Files
	files, err := listFiles(workingDir)

	if err != nil {
		return fmt.Errorf("Failed to list files: %v", err)
	}

	for _, f := range files {
		fileContent, err := os.ReadFile(f)

		if err != nil {
			return fmt.Errorf("Failed processing file '%s': %v", f, err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "%s---\n", string(fileContent))

	}

	return nil

}

func listFiles(root string) ([]string, error) {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".yml" || filepath.Ext(path) == ".json" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return files, err
	}

	return files, nil
}
