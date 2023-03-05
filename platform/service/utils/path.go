package utils

import (
	"github.com/obgnail/plugin-platform/common/config"
	"github.com/obgnail/plugin-platform/common/utils/file_path"
	"github.com/obgnail/plugin-platform/platform/service/types"
)

var (
	PluginFileDir      = config.StringOrPanic("platform.plugins_storage_dir")
	PluginWorkspaceDir = config.StringOrPanic("platform.plugin_runtime_dir")
)

func IsPluginYamlPath(yamlPath string) bool {
	p := file_path.JoinPath(types.PluginFileConfigDir, types.PluginFileConfigYaml)
	return yamlPath == p
}

func GetPluginUpgradeFilePath(appUUID string) string {
	return file_path.JoinPath(PluginFileDir, appUUID, types.PluginUpgradeDir)
}

func GetPluginDir(appUUID string, version string) string {
	return file_path.JoinPath(PluginFileDir, appUUID, version)
}

func GetPluginConfigPath(appUUID string, version string) string {
	return file_path.JoinPath(PluginFileDir, appUUID, version, types.PluginFileConfigDir, types.PluginFileConfigYaml)
}

func GetPluginSoFile(appUUID string, version string) string {
	return file_path.JoinPath(PluginFileDir, appUUID, version, types.PluginFileBackDir, types.PluginServerDir, types.PluginSoFile)
}

func GetPluginWorkspace(appUUID string, version string) string {
	return file_path.JoinPath(PluginWorkspaceDir, appUUID, version)
}
