package cmd

import (
	"fmt"
	"github.com/geowa4/ocm-workon/pkg/cluster"
	"github.com/geowa4/ocm-workon/pkg/config"
	"github.com/geowa4/ocm-workon/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"os"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes directories and files for other commands",
	Long:  `Creates the base directory for working on clusters and common zshrc file. Additionally persists all passed in configuration to the config file. All flags for this command match the cluster command. Existing values in the configuration file are preserved.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		bindViperToClusterFlags(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		cobra.CheckErr(cluster.Initialize(viper.GetString("cluster_base_directory")))
		writeAllSettings()
		validateOcmConfigs()
		validateBackplaneConfigs()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	addConfigurationFlagsToClusterCmd(initCmd)
}

func writeAllSettings() {
	allSettings := viper.AllSettings()
	marshaledSettings, err := yaml.Marshal(allSettings)
	cobra.CheckErr(err)
	configFile, err := os.OpenFile(config.GetOcmConfigDir()+utils.PathSep+ConfigFileName, os.O_WRONLY, 0644)
	cobra.CheckErr(err)
	_, err = configFile.Write(marshaledSettings)
	cobra.CheckErr(err)
}

func validateOcmConfigs() {
	ocmConfigDir := config.GetOcmConfigDir()
	prodConfig := ocmConfigDir + utils.PathSep + "ocm." + cluster.ProductionEnvironment + ".json"
	if _, err := os.Lstat(prodConfig); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "WARN: please create %s\n", prodConfig)
	}

	stageConfig := ocmConfigDir + utils.PathSep + "ocm." + cluster.StagingEnvironment + ".json"
	if _, err := os.Lstat(prodConfig); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "WARN: please create %s\n", stageConfig)
	}
}

func validateBackplaneConfigs() {
	backplaneConfigDir := config.GetBackplaneConfigDir()
	prodConfig := backplaneConfigDir + utils.PathSep + "config." + cluster.ProductionEnvironment + ".json"
	if _, err := os.Lstat(prodConfig); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "WARN: please create %s\n", prodConfig)
	}

	stageConfig := backplaneConfigDir + utils.PathSep + "config." + cluster.StagingEnvironment + ".json"
	if _, err := os.Lstat(prodConfig); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "WARN: please create %s\n", stageConfig)
	}
}
