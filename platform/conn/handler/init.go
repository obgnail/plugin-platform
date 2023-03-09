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

func CallPluginFunction(instanceID string, abilityID, abilityType, abilityFuncKey string, arg []byte) chan *common_type.AbilityResponse {
	return ph.CallPluginFunction(instanceID, abilityID, abilityType, abilityFuncKey, arg)
}

func CallPluginHTTP(instanceID string, req *common_type.HttpRequest, internal bool, abilityFunc string) chan *common_type.HttpResponse {
	return ph.CallPluginHTTP(instanceID, req, internal, abilityFunc)
}

func CallPluginInternalHTTP(instanceID string, req *common_type.HttpRequest, abilityFunc string) chan *common_type.HttpResponse {
	return ph.CallPluginHTTP(instanceID, req, true, abilityFunc)
}

func CallPluginExternalHTTP(instanceID string, req *common_type.HttpRequest) chan *common_type.HttpResponse {
	return ph.CallPluginHTTP(instanceID, req, false, "")
}

func CallPluginEvent(instanceID string, eventType string, payload []byte) chan common_type.PluginError {
	return ph.CallPluginEvent(instanceID, eventType, payload, false)
}

func ForceCallPluginEvent(instanceID string, eventType string, payload []byte) chan common_type.PluginError {
	return ph.CallPluginEvent(instanceID, eventType, payload, true)
}

func CallPluginConfigChanged(instanceID string, configKey string, originValue, newValue []string) chan common_type.PluginError {
	return ph.CallPluginConfigChanged(instanceID, configKey, originValue, newValue)
}

func SendToAlivePlugin(instanceID string, messageBuilder MessageBuilder) (*protocol.PlatformMessage, common_type.PluginError) {
	return ph.SendToAlivePlugin(instanceID, messageBuilder)
}

func GetHost(instanceID string) common_type.IHost {
	return ph.GetHost(instanceID)
}

func GetHostBoot(hostBootID string) common_type.IHostBoot {
	return ph.GetHostBoot(hostBootID)
}

func GetHosts() []common_type.IHost {
	return ph.GetAllHost()
}

func GetHostBoots() []common_type.IHostBoot {
	return ph.GetAllHostBoot()
}

func GetAlivePlugins() map[string]common_type.IInstanceDescription {
	return ph.GetAllAlivePlugin()
}

func GetSupportPlugins() map[string]common_type.IInstanceDescription {
	return ph.GetAllSupportPlugin()
}

func KillPlugin(instanceID string) chan common_type.PluginError {
	return ph.KillPlugin(instanceID)
}

func KillHost(hostID string) {
	ph.KillHost(hostID)
}
