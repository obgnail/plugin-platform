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

func InstallPlugin(appID, instanceID, name, lang, langVer, appVer string) chan common_type.PluginError {
	return ph.InstallPlugin(appID, instanceID, name, lang, langVer, appVer)
}

func UpgradePlugin(appID, instanceID, name, lang, langVer, appVer string, oldVersion *protocol.PluginDescriptor) chan common_type.PluginError {
	return ph.UpgradePlugin(appID, instanceID, name, lang, langVer, appVer, oldVersion)
}

func EnablePlugin(instanceID string) chan common_type.PluginError {
	return ph.EnablePlugin(instanceID)
}

func DisablePlugin(instanceID string) chan common_type.PluginError {
	return ph.DisablePlugin(instanceID)
}

func UnInstallPlugin(instanceID string) chan common_type.PluginError {
	return ph.UnInstallPlugin(instanceID)
}

func CheckStatePlugin(instanceID string) chan common_type.PluginError {
	return ph.CheckStatePlugin(instanceID)
}

func CheckCompatibilityPlugin(instanceID string) chan common_type.PluginError {
	return ph.CheckCompatibilityPlugin(instanceID)
}

func CallPluginHttp(instanceID string, req *common_type.HttpRequest, abilityFunc string) chan *common_type.HttpResponse {
	return ph.CallPluginHttp(instanceID, req, abilityFunc)
}

func CallPluginEvent(instanceID string, eventType string, payload []byte) chan common_type.PluginError {
	return ph.CallPluginEvent(instanceID, eventType, payload)
}

func KillPlugin(instanceID string) { ph.KillPlugin(instanceID) }

func KillHost(hostID string) { ph.KillHost(hostID) }

func GetHost(instanceID string) common_type.IHost { return ph.GetHost(instanceID) }

func GetHostBoot(hostBootID string) common_type.IHostBoot { return ph.GetHostBoot(hostBootID) }

func GetHosts() []common_type.IHost { return ph.GetAllHost() }

func GetHostBoots() []common_type.IHostBoot { return ph.GetAllHostBoot() }

func GetAlivePlugins() map[string]common_type.IInstanceDescription { return ph.GetAllAlivePlugin() }

func GetSupportPlugins() map[string]common_type.IInstanceDescription { return ph.GetAllSupportPlugin() }
