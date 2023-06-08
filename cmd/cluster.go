package cmd

import (
	"fmt"
	"github.com/geowa4/ocm-workon/pkg/cluster"
	"github.com/geowa4/ocm-workon/pkg/shell"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// clusterCmd represents the cluster command
var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Sets up a directory to work on a cluster",
	Long: `Sets up a directory with its own kube config and environment variables commonly used in troubleshooting scripts.

Example: cluster --production 15d716b7-b933-41ef-924c-53c2b59afe4f`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		environment := cluster.StagingEnvironment
		if viper.GetBool("cluster_production") {
			environment = cluster.ProductionEnvironment
		}

		ncd, err := cluster.NewNormalizedCluster(args[0])
		if err != nil {
			fmt.Printf("error retrieving cluster data: %q\n", err)
			os.Exit(1)
		}
		baseDir := viper.GetString("cluster_base_directory")
		w := &cluster.WorkConfig{
			Environment: environment,
			ClusterData: ncd,
			ClusterBase: baseDir,
			UseDirenv:   viper.GetBool("cluster_use_direnv"),
			UseAsdf:     viper.GetBool("cluster_use_asdf"),
		}
		if clusterDir, err := w.Build(); err != nil {
			fmt.Println(err)
			os.Exit(2)
		} else {
			recordedCluster := cluster.NewRecordedCluster(environment, ncd)
			if err = recordedCluster.RecordAccess(baseDir); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "[WARN] access of this cluster will not be recorded due to error: %q\n", err)
			}
			if err = shell.Exec(baseDir, clusterDir); err != nil {
				fmt.Println(err)
				os.Exit(3)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(clusterCmd)

	clusterCmd.Flags().BoolP("production", "p", false, "Whether this is a staging (default) or production cluster")
	_ = viper.BindPFlag("cluster_production", clusterCmd.Flags().Lookup("production"))

	clusterCmd.Flags().StringP("base-dir", "b", "", "The base directory for creating working environments for each cluster")
	_ = viper.BindPFlag("cluster_base_directory", clusterCmd.Flags().Lookup("base-dir"))
	_ = clusterCmd.MarkFlagDirname("base-dir")

	clusterCmd.Flags().BoolP("use-direnv", "d", true, "Whether to use direnv")
	_ = viper.BindPFlag("cluster_use_direnv", clusterCmd.Flags().Lookup("use-direnv"))

	clusterCmd.Flags().BoolP("use-asdf", "a", false, "Whether to use asdf in the .envrc for direnv")
	_ = viper.BindPFlag("cluster_use_asdf", clusterCmd.Flags().Lookup("use-asdf"))

	clusterCmd.Flags().StringP("shell", "s", "/bin/zsh", "The shell or other executable to run when the directory is built")
	_ = viper.BindPFlag("cluster_shell", clusterCmd.Flags().Lookup("shell"))
	_ = clusterCmd.MarkFlagFilename("shell")

	clusterCmd.Flags().StringArray("shell-args", []string{"--login", "-i"}, "Arguments to pass to the shell")
	_ = viper.BindPFlag("cluster_shell_args", clusterCmd.Flags().Lookup("shell-args"))
}
