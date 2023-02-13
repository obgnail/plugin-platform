package utils

import (
	"github.com/obgnail/plugin-platform/platform/config"
	"github.com/obgnail/plugin-platform/platform/service/types"
	"strings"
)

func JoinPath(paths ...string) string {
	return strings.Join(paths, "/")
}

func IsPluginYamlPath(yamlPath string) bool {
	p := JoinPath(types.PluginFileConfigDir, types.PluginFileConfigYaml)
	return yamlPath == p
}

func GetPluginUpgradeFilePath(appUUID string) string {
	PluginFileDir := config.StringOrPanic("plugin_file_directory")
	return JoinPath(PluginFileDir, appUUID, types.PluginUpgradeDir)
}

func GetPluginDir(appUUID string, version string) string {
	PluginFileDir := config.StringOrPanic("plugin_file_directory")
	return JoinPath(PluginFileDir, appUUID, version)
}

func GetPluginConfigPath(appUUID string, version string) string {
	PluginFileDir := config.StringOrPanic("plugin_file_directory")
	return JoinPath(PluginFileDir, appUUID, version, types.PluginFileConfigDir, types.PluginFileConfigYaml)
}
