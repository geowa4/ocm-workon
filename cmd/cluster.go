package cmd

import (
	"fmt"
	"github.com/geowa4/ocm-workon/pkg/cluster"
	"github.com/geowa4/ocm-workon/pkg/utils"
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
	PreRun: func(cmd *cobra.Command, args []string) {
		bindViperToClusterFlags(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		environment := cluster.ProductionEnvironment
		if viper.GetBool("cluster_staging") {
			environment = cluster.StagingEnvironment
		}

		ncd, err := cluster.NewNormalizedCluster(args[0])
		cobra.CheckErr(err)
		baseDir := viper.GetString("cluster_base_directory")
		w := &cluster.WorkConfig{
			Environment: environment,
			ClusterData: ncd,
			ClusterBase: baseDir,
			UseDirenv:   viper.GetBool("cluster_use_direnv"),
			UseAsdf:     viper.GetBool("cluster_use_asdf"),
		}
		clusterDir, err := w.Build()
		cobra.CheckErr(err)

		recordedCluster := cluster.NewRecordedCluster(environment, ncd)
		if err = recordedCluster.RecordAccess(baseDir); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Warning: access of this cluster will not be recorded due to error: %q\n", err)
		}
		cobra.CheckErr(utils.ShellExec(baseDir, clusterDir))
	},
}

func init() {
	rootCmd.AddCommand(clusterCmd)
	addConfigurationFlagsToClusterCmd(clusterCmd)

	clusterCmd.Flags().Bool("stage", false, "Whether this is a stage or production (default) cluster")
}

// Flags that can be backed up by a persistent config file or environment variables
func addConfigurationFlagsToClusterCmd(cmd *cobra.Command) {
	homeDir, _ := os.UserHomeDir()
	cmd.Flags().StringP("base-dir", "b", homeDir+utils.PathSep+"Clusters", "The base directory for creating working environments for each cluster")
	_ = cmd.MarkFlagDirname("base-dir")

	cmd.Flags().BoolP("use-direnv", "d", true, "Whether to use direnv")

	cmd.Flags().BoolP("use-asdf", "a", false, "Whether to use asdf in the .envrc for direnv")

	cmd.Flags().StringP("shell", "s", "/bin/zsh", "The shell or other executable to run when the directory is built")
	_ = cmd.MarkFlagFilename("shell")

	cmd.Flags().StringArray("shell-args", []string{"--login", "-i"}, "Arguments to pass to the shell")
}

func bindViperToClusterFlags(cmd *cobra.Command) {
	_ = viper.BindPFlag("cluster_base_directory", cmd.Flags().Lookup("base-dir"))
	_ = viper.BindPFlag("cluster_use_direnv", cmd.Flags().Lookup("use-direnv"))
	_ = viper.BindPFlag("cluster_use_asdf", cmd.Flags().Lookup("use-asdf"))
	_ = viper.BindPFlag("cluster_shell", cmd.Flags().Lookup("shell"))
	_ = viper.BindPFlag("cluster_shell_args", cmd.Flags().Lookup("shell-args"))
}
