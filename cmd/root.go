package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ocm-workon",
	Short: "Work on a cluster",
	Long: `Sets up a directory with its own kube config and environment variables commonly used in troubleshooting scripts.

Example: cluster --production 15d716b7-b933-41ef-924c-53c2b59afe4f`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
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
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name "workon" (without extension).
		configDir := home + string(os.PathSeparator) + ".config"
		if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
			configDir = xdgConfigHome
		}
		ocmConfigDir := configDir + string(os.PathSeparator) + "ocm"
		viper.AddConfigPath(ocmConfigDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("workon.yaml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error using config file:", viper.ConfigFileUsed())
		os.Exit(1)
	}
}
