package cmd

import (
	"github.com/charmbracelet/log"
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

func validateBackplaneConfigs() {
	prodConfig := config.GetBackplaneConfigFile(cluster.ProductionEnvironment)
	if _, err := os.Lstat(prodConfig); err != nil {
		log.Warnf("please create %s", prodConfig)
	}

	stageConfig := config.GetBackplaneConfigFile(cluster.ProductionEnvironment)
	if _, err := os.Lstat(stageConfig); err != nil {
		log.Warnf("please create %s", stageConfig)
	}
}
