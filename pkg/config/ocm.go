package config

import (
	"github.com/geowa4/ocm-workon/pkg/utils"
	"github.com/spf13/cobra"
	"os"
)

func GetOcmConfigDir() string {
	return getHomeConfigDir() + utils.PathSep + "ocm"
}

func GetBackplaneConfigDir() string {
	return getHomeConfigDir() + utils.PathSep + "backplane"
}

func GetBackplaneConfigFile(environment string) string {
	return GetBackplaneConfigDir() + utils.PathSep + "config." + environment + ".json"
}

func getHomeConfigDir() string {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	configDir := home + string(os.PathSeparator) + ".config"
	if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
		configDir = xdgConfigHome
	}
	return configDir
}
