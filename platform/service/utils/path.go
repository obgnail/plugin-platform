package utils

import (
	"github.com/obgnail/plugin-platform/common/config"
	"github.com/obgnail/plugin-platform/common/file_utils"
	"github.com/obgnail/plugin-platform/platform/service/types"
)

var PluginFileDir = config.StringOrPanic("platform.plugins_storage_dir")

func IsPluginYamlPath(yamlPath string) bool {
	p := file_utils.JoinPath(types.PluginFileConfigDir, types.PluginFileConfigYaml)
	return yamlPath == p
}

func GetPluginUpgradeFilePath(appUUID string) string {
	return file_utils.JoinPath(PluginFileDir, appUUID, types.PluginUpgradeDir)
}

func GetPluginDir(appUUID string, version string) string {
	return file_utils.JoinPath(PluginFileDir, appUUID, version)
}

func GetPluginConfigPath(appUUID string, version string) string {
	return file_utils.JoinPath(PluginFileDir, appUUID, version, types.PluginFileConfigDir, types.PluginFileConfigYaml)
}

func GetPluginSoFile(appUUID string, version string) string {
	return file_utils.JoinPath(PluginFileDir, appUUID, version, types.PluginFileBackDir, types.PluginServerDir, types.PluginSoFile)
}
