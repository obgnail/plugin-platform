package utils

import (
	"github.com/obgnail/plugin-platform/common/config"
	"github.com/obgnail/plugin-platform/common/file_utils"
	"github.com/obgnail/plugin-platform/platform/service/types"
)

func IsPluginYamlPath(yamlPath string) bool {
	p := file_utils.JoinPath(types.PluginFileConfigDir, types.PluginFileConfigYaml)
	return yamlPath == p
}

func GetPluginUpgradeFilePath(appUUID string) string {
	PluginFileDir := config.StringOrPanic("platform.plugins_storage_dir")
	return file_utils.JoinPath(PluginFileDir, appUUID, types.PluginUpgradeDir)
}

func GetPluginDir(appUUID string, version string) string {
	PluginFileDir := config.StringOrPanic("platform.plugins_storage_dir")
	return file_utils.JoinPath(PluginFileDir, appUUID, version)
}

func GetPluginConfigPath(appUUID string, version string) string {
	PluginFileDir := config.StringOrPanic("platform.plugins_storage_dir")
	return file_utils.JoinPath(PluginFileDir, appUUID, version, types.PluginFileConfigDir, types.PluginFileConfigYaml)
}
