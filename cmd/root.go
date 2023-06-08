package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const ConfigFileName = "workon.yaml"

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ocm-workon",
	Short: "Work on a cluster",

	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ${$XDG_CONFIG_DIR:-$HOME/.config}/ocm/workon.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		ocmConfigDir := getOcmConfigDir()
		viper.AddConfigPath(ocmConfigDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName(ConfigFileName)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			cobra.CheckErr(fmt.Errorf("bad config file (%s): %q", viper.ConfigFileUsed(), err))
		}
	}
}

func getOcmConfigDir() string {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	configDir := home + string(os.PathSeparator) + ".config"
	if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
		configDir = xdgConfigHome
	}
	ocmConfigDir := configDir + string(os.PathSeparator) + "ocm"
	return ocmConfigDir
}
