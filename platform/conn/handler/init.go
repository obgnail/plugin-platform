package handler

import (
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/protocol"
)

var ph *PlatformHandler

func InitPlatformHandler() error {
	ph = Default()
	ph.Run()
	return nil
}

// 生命周期函数全部改成同步,因为业务上用不到异步

func EnablePlugin(appID, instanceID, name, lang, langVer, appVer string) common_type.PluginError {
	return <-ph.EnablePlugin(make(chan common_type.PluginError, 1), appID, instanceID, name, lang, langVer, appVer)
}

func DisablePlugin(appID, instanceID, name, lang, langVer, appVer string) common_type.PluginError {
	return <-ph.DisablePlugin(make(chan common_type.PluginError, 1), appID, instanceID, name, lang, langVer, appVer)
}

func InstallPlugin(appID, instanceID, name, lang, langVer, appVer string) common_type.PluginError {
	return <-ph.InstallPlugin(make(chan common_type.PluginError, 1), appID, instanceID, name, lang, langVer, appVer)
}

func UnInstallPlugin(appID, instanceID, name, lang, langVer, appVer string) common_type.PluginError {
	return <-ph.UnInstallPlugin(make(chan common_type.PluginError, 1), appID, instanceID, name, lang, langVer, appVer)
}

func UpgradePlugin(appID, instanceID, name, lang, langVer, appVer string, oldVersion *protocol.PluginDescriptor) common_type.PluginError {
	return <-ph.UpgradePlugin(make(chan common_type.PluginError, 1), appID, instanceID, name, lang, langVer, appVer, oldVersion)
}

func CheckStatePlugin(appID, instanceID, name, lang, langVer, appVer string) common_type.PluginError {
	return <-ph.CheckStatePlugin(make(chan common_type.PluginError, 1), appID, instanceID, name, lang, langVer, appVer)
}

func CheckCompatibilityPlugin(appID, instanceID, name, lang, langVer, appVer string) common_type.PluginError {
	return <-ph.CheckCompatibilityPlugin(make(chan common_type.PluginError, 1), appID, instanceID, name, lang, langVer, appVer)
}

func KillPlugin(instanceID string) { ph.KillPlugin(instanceID) }

func KillHost(hostID string) { ph.KillHost(hostID) }

func GetHost(instanceID string) common_type.IHost { return ph.GetHost(instanceID) }

func GetHostBoot(hostBootID string) common_type.IHostBoot { return ph.GetHostBoot(hostBootID) }

func GetHosts() []common_type.IHost { return ph.GetAllHost() }

func GetHostBoots() []common_type.IHostBoot { return ph.GetAllHostBoot() }

func GetAlivePlugins() map[string]common_type.IInstanceDescription { return ph.GetAllAlivePlugin() }

func GetSupportPlugins() map[string]common_type.IInstanceDescription { return ph.GetAllSupportPlugin() }
