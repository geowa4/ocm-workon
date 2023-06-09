package shell

import (
	"github.com/spf13/viper"
	"os"
	"syscall"
)

func Exec(baseDir string, clusterDir string) error {
	environ := append(os.Environ(),
		"ZDOTDIR="+viper.GetString("cluster_base_directory"),
		"CLUSTER_HOME="+clusterDir,
		"CLUSTER_BASE_DIRECTORY="+baseDir)
	return syscall.Exec(
		viper.GetString("cluster_shell"),
		viper.GetStringSlice("cluster_shell_args"),
		environ,
	)
}
