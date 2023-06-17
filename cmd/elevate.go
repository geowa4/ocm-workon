package cmd

import (
	"fmt"
	"github.com/geowa4/ocm-workon/pkg/cluster"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// elevationCmd represents the compliance command
var elevationCmd = &cobra.Command{
	Use:   "elevate",
	Short: "Records the elevation",
	Long: `Records the elevation along with useful reminders to provide adequate justification. 

Example: elevate --source 'https://issues.redhat.com/browse/OCPBUGS-7158' --reason 'a very good reason'`,
	PreRun: func(cmd *cobra.Command, args []string) {
		_ = viper.BindPFlag("cluster_id", cmd.Flags().Lookup("cluster-id"))
		_ = viper.BindPFlag("cluster_base_directory", cmd.Flags().Lookup("cluster-base-directory"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		source := cmd.Flags().Lookup("source").Value.String()
		reason := cmd.Flags().Lookup("reason").Value.String()
		cobra.CheckErr(cluster.RecordElevation(
			viper.GetString("cluster_base_directory"),
			viper.GetString("cluster_id"),
			source,
			reason,
		))
		fmt.Println(fmt.Sprintf(
			"ELEVATION_REASON='%s %s'",
			source,
			reason,
		))
	},
}

func init() {
	rootCmd.AddCommand(elevationCmd)

	elevationCmd.Flags().StringP("cluster-id", "c", "", "The ID of the cluster where you need to have elevated permissions. Uses $CLUSTER_ID by default.")

	elevationCmd.Flags().StringP("cluster-base-directory", "b", "", "The base directory for all clusters worked. Uses $CLUSTER_BASE_DIRECTORY by default.")

	elevationCmd.Flags().StringP("source", "s", "", "A link or other source that triggered the work.")
	elevationCmd.Flags().StringP("reason", "r", "", "A brief description of why you needed to elevate.")
	elevationCmd.MarkFlagsRequiredTogether("source", "reason")
	_ = elevationCmd.MarkFlagRequired("source")
	_ = elevationCmd.MarkFlagRequired("reason")
}
