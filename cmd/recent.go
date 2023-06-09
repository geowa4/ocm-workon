package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/geowa4/ocm-workon/pkg/cluster"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "recent",
	Short: "List clusters you have worked on recently",
	Long:  `List clusters as JSON that you have worked on in the past _two weeks_.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		_ = viper.BindPFlag("cluster_base_directory", cmd.Flags().Lookup("base-dir"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		clusters, err := cluster.FindClustersUpdatedSinceTwoWeeksAgo(viper.GetString("cluster_base_directory"))
		marshaled, err := json.MarshalIndent(clusters, "", "  ")
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println(string(marshaled))
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringP("base-dir", "b", "", "The base directory for creating working environments for each cluster")
	_ = listCmd.MarkFlagDirname("base-dir")
}
