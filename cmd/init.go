package cmd

import (
	"fmt"
	"github.com/geowa4/ocm-workon/pkg/cluster"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes directories and files for other commands",
	Long: `Creates the base directory for working on clusters and common zshrc file.

Example: init --base-dir ~/Clusters`,
	Run: func(cmd *cobra.Command, args []string) {
		baseDir := viper.GetString("cluster_base_directory")
		if err := cluster.Initialize(baseDir); err != nil {
			fmt.Printf("error creating base directory: %q\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringP("base-dir", "b", "", "The base directory for creating working environments for each cluster")
	_ = viper.BindPFlag("cluster_base_directory", initCmd.Flags().Lookup("base-dir"))
	_ = initCmd.MarkFlagDirname("base-dir")
}
