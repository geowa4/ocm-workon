package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/geowa4/ocm-workon/pkg/cluster"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "recent",
	Short: "List clusters you have worked on recently",
	Long: `List clusters as JSON that you have worked on in the the last three days by default.

Example: recent --since 2w`,
	PreRun: func(cmd *cobra.Command, args []string) {
		_ = viper.BindPFlag("cluster_base_directory", cmd.Flags().Lookup("base-dir"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		var (
			resources any
			err       error
		)
		if cmd.Flags().Lookup("elevations").Value.String() == "true" {
			resources, err = cluster.FindElevationsSince(
				viper.GetString("cluster_base_directory"),
				cmd.Flags().Lookup("since").Value.String(),
			)
		} else {
			resources, err = cluster.FindRecordedClustersSince(
				viper.GetString("cluster_base_directory"),
				cmd.Flags().Lookup("since").Value.String(),
			)
		}
		cobra.CheckErr(err)
		marshaled, err := json.Marshal(resources)
		cobra.CheckErr(err)
		fmt.Println(string(marshaled))
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringP("base-dir", "b", "", "The base directory for creating working environments for each cluster")
	_ = listCmd.MarkFlagDirname("base-dir")

	listCmd.Flags().String("since", "24h", "How far back to search for clusters or elevations")

	listCmd.Flags().Bool("elevations", false, "Whether to load elevations or the default clusters")
}
