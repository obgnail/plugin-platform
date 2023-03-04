package handler

import (
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/protocol"
)

var platformHandler *PlatformHandler

func InitPlatformHandler() error {
	platformHandler = Default()
	platformHandler.Run()
	return nil
}

// 生命周期函数全部改成同步,因为业务上用不到异步

func EnablePlugin(appID, instanceID, name, lang, langVer, appVer string) common_type.PluginError {
	return <-platformHandler.EnablePlugin(make(chan common_type.PluginError, 1), appID, instanceID, name, lang, langVer, appVer)
}

func DisablePlugin(appID, instanceID, name, lang, langVer, appVer string) common_type.PluginError {
	return <-platformHandler.DisablePlugin(make(chan common_type.PluginError, 1), appID, instanceID, name, lang, langVer, appVer)
}

func StartPlugin(appID, instanceID, name, lang, langVer, appVer string) common_type.PluginError {
	return <-platformHandler.StartPlugin(make(chan common_type.PluginError, 1), appID, instanceID, name, lang, langVer, appVer)
}

func StopPlugin(appID, instanceID, name, lang, langVer, appVer string) common_type.PluginError {
	return <-platformHandler.StopPlugin(make(chan common_type.PluginError, 1), appID, instanceID, name, lang, langVer, appVer)
}

func InstallPlugin(appID, instanceID, name, lang, langVer, appVer string) common_type.PluginError {
	return <-platformHandler.InstallPlugin(make(chan common_type.PluginError, 1), appID, instanceID, name, lang, langVer, appVer)
}

func UnInstallPlugin(appID, instanceID, name, lang, langVer, appVer string) common_type.PluginError {
	return <-platformHandler.UnInstallPlugin(make(chan common_type.PluginError, 1), appID, instanceID, name, lang, langVer, appVer)
}

func UpgradePlugin(appID, instanceID, name, lang, langVer, appVer string, oldVersion *protocol.PluginDescriptor) common_type.PluginError {
	return <-platformHandler.UpgradePlugin(make(chan common_type.PluginError, 1), appID, instanceID, name, lang, langVer, appVer, oldVersion)
}

func CheckStatePlugin(appID, instanceID, name, lang, langVer, appVer string) common_type.PluginError {
	return <-platformHandler.CheckStatePlugin(make(chan common_type.PluginError, 1), appID, instanceID, name, lang, langVer, appVer)
}

func CheckCompatibilityPlugin(appID, instanceID, name, lang, langVer, appVer string) common_type.PluginError {
	return <-platformHandler.CheckCompatibilityPlugin(make(chan common_type.PluginError, 1), appID, instanceID, name, lang, langVer, appVer)
}

func KillHost(hostID string) { platformHandler.KillHost(hostID) }

func KillPlugin(instanceID string) { platformHandler.KillPlugin(instanceID) }

func GetAllHost() []common_type.IHost { return platformHandler.GetAllHost() }

func GetAllHostBoot() []common_type.IHostBoot { return platformHandler.GetAllHostBoot() }

func GetAllAlivePlugin() map[string]common_type.IInstanceDescription {
	return platformHandler.GetAllAlivePlugin()
}

func GetAllSupportPlugin() map[string]common_type.IInstanceDescription {
	return platformHandler.GetAllSupportPlugin()
}
