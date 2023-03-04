package handler

import (
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/protocol"
)

var platformHandler *Handler

func InitPlatformHandler() error {
	platformHandler = Default()
	platformHandler.Run()
	return nil
}

func EnablePlugin(appID, instanceID, name, lang, langVer, appVer string) {
	platformHandler.lifeCycleReq(protocol.ControlMessage_Enable, appID, instanceID, name, lang, langVer, appVer, nil)
}

func DisablePlugin(appID, instanceID, name, lang, langVer, appVer string) {
	platformHandler.lifeCycleReq(protocol.ControlMessage_Disable, appID, instanceID, name, lang, langVer, appVer, nil)
}

func StartPlugin(appID, instanceID, name, lang, langVer, appVer string) {
	platformHandler.lifeCycleReq(protocol.ControlMessage_Start, appID, instanceID, name, lang, langVer, appVer, nil)
}

func StopPlugin(appID, instanceID, name, lang, langVer, appVer string) {
	platformHandler.lifeCycleReq(protocol.ControlMessage_Stop, appID, instanceID, name, lang, langVer, appVer, nil)
}

func InstallPlugin(appID, instanceID, name, lang, langVer, appVer string) {
	platformHandler.lifeCycleReq(protocol.ControlMessage_Install, appID, instanceID, name, lang, langVer, appVer, nil)
}

func UnInstallPlugin(appID, instanceID, name, lang, langVer, appVer string) {
	platformHandler.lifeCycleReq(protocol.ControlMessage_UnInstall, appID, instanceID, name, lang, langVer, appVer, nil)
}

func UpgradePlugin(appID, instanceID, name, lang, langVer, appVer string, oldVersion *protocol.PluginDescriptor) {
	platformHandler.lifeCycleReq(protocol.ControlMessage_Upgrade, appID, instanceID, name, lang, langVer, appVer, oldVersion)
}

func CheckStatePlugin(appID, instanceID, name, lang, langVer, appVer string) {
	platformHandler.lifeCycleReq(protocol.ControlMessage_CheckState, appID, instanceID, name, lang, langVer, appVer, nil)
}

func CheckCompatibilityPlugin(appID, instanceID, name, lang, langVer, appVer string) {
	platformHandler.lifeCycleReq(protocol.ControlMessage_CheckCompatibility, appID, instanceID, name, lang, langVer, appVer, nil)
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
