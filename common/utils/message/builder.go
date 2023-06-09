package message

import (
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/config"
	"github.com/obgnail/plugin-platform/common/connect"
	"github.com/obgnail/plugin-platform/common/protocol"
	"github.com/obgnail/plugin-platform/common/utils/math"
	"time"
)

var (
	platformID   = config.StringOrPanic("platform.id")
	platformName = config.StringOrPanic("platform.name")
)

func GetInitMessage(source, distinct *protocol.RouterNode) *protocol.PlatformMessage {
	message := &protocol.PlatformMessage{
		Header: &protocol.RouterMessage{
			SeqNo:    math.CreateCaptcha(),
			Source:   source,
			Distinct: distinct,
		},
		Control: &protocol.ControlMessage{},
	}
	return message
}

func GetHostBootInfo(bootID, bootName string) *protocol.RouterNode {
	return &protocol.RouterNode{
		ID: bootID,
		Tags: map[string]string{
			"role": connect.RoleHostBoot,
			"id":   bootID,
			"name": bootName,
		},
	}
}

func GetHostInfo(hostID, hostName string) *protocol.RouterNode {
	return &protocol.RouterNode{
		ID: hostID,
		Tags: map[string]string{
			"role": connect.RoleHost,
			"id":   hostID,
			"name": hostName,
		},
	}
}

func GetPlatformInfo() *protocol.RouterNode {
	return &protocol.RouterNode{
		ID: platformID,
		Tags: map[string]string{
			"role": connect.RolePlatform,
			"id":   platformID,
			"name": platformName,
		},
	}
}

// platform -> host
func BuildP2HDefaultMessage(hostID, hostName string) *protocol.PlatformMessage {
	p := GetPlatformInfo()
	h := GetHostInfo(hostID, hostName)
	msg := GetInitMessage(p, h)
	return msg
}

// host -> platform
func BuildH2PDefaultMessage(hostID, hostName string) *protocol.PlatformMessage {
	p := GetPlatformInfo()
	h := GetHostInfo(hostID, hostName)
	msg := GetInitMessage(h, p)
	return msg
}

// platform -> boot
func BuildP2BDefaultMessage(bootID, bootName string) *protocol.PlatformMessage {
	p := GetPlatformInfo()
	b := GetHostBootInfo(bootID, bootName)
	msg := GetInitMessage(p, b)
	return msg
}

// boot -> platform
func BuildB2PDefaultMessage(bootID, bootName string) *protocol.PlatformMessage {
	p := GetPlatformInfo()
	b := GetHostBootInfo(bootID, bootName)
	msg := GetInitMessage(b, p)
	return msg
}

// SwapMessageHeader 消息头 来源/去向 调换
func SwapMessageHeader(newMsg, comeMsg *protocol.PlatformMessage) {
	newMsg.Header.Source.ID = comeMsg.Header.Distinct.ID
	newMsg.Header.Source.Tags = comeMsg.Header.Distinct.Tags
	newMsg.Header.Distinct.ID = comeMsg.Header.Source.ID
	newMsg.Header.Distinct.Tags = comeMsg.Header.Source.Tags
	newMsg.Header.SeqNo = comeMsg.Header.SeqNo
}

func BuildErrorMessage(err common_type.PluginError) *protocol.ErrorMessage {
	if err == nil {
		return nil
	}

	return &protocol.ErrorMessage{
		Code:  int64(err.Code()),
		Error: err.Error(),
		Msg:   err.Msg(),
	}
}

func BuildInstanceDescriptor(description common_type.IInstanceDescription, hostID string) *protocol.PluginInstanceDescriptor {
	desc := description.PluginDescription()
	return &protocol.PluginInstanceDescriptor{
		Application: &protocol.PluginDescriptor{
			ApplicationID:      desc.ApplicationID(),
			Name:               desc.Name(),
			Language:           desc.Language(),
			LanguageVersion:    NewProtocolVersion(desc.LanguageVersion()),
			ApplicationVersion: NewProtocolVersion(desc.ApplicationVersion()),
			HostVersion:        NewProtocolVersion(desc.HostVersion()),
			MinSystemVersion:   NewProtocolVersion(desc.MinSystemVersion()),
		},
		InstanceID: description.InstanceID(),
		HostID:     hostID,
	}
}

func BuildCallPluginHTTPMessage(
	req *common_type.HttpRequest,
	internal bool,
	host common_type.HostInfo,
	target common_type.IInstanceDescription,
	abilityFunc string,
) *protocol.PlatformMessage {
	msg := BuildP2HDefaultMessage(host.ID, host.Name)
	msg.Plugin = &protocol.PluginMessage{
		Target: BuildInstanceDescriptor(target, host.ID),
		Http: &protocol.HttpContextMessage{
			Request: &protocol.HttpRequestMessage{
				Method:      req.Method,
				Url:         req.Url,
				Headers:     nil,
				Body:        req.Body,
				Internal:    internal,
				AbilityFunc: abilityFunc,
			},
		},
	}
	return msg
}

func BuildCallPluginFunctionMessage(
	abilityID, abilityType, abilityFuncKey string, arg []byte,
	host common_type.HostInfo,
	target common_type.IInstanceDescription,
) *protocol.PlatformMessage {
	msg := BuildP2HDefaultMessage(host.ID, host.Name)
	msg.Plugin = &protocol.PluginMessage{
		Target: BuildInstanceDescriptor(target, host.ID),
		Ability: &protocol.StandardAbilityMessage{
			AbilityRequest: &protocol.StandardAbilityMessage_AbilityRequestMessage{
				Id:      abilityID,
				Type:    abilityType,
				FuncKey: abilityFuncKey,
				Args:    arg,
			},
		},
	}
	return msg
}

func BuildCallPluginEventMessage(
	eventType string, payload []byte,
	host common_type.HostInfo,
	target common_type.IInstanceDescription,
) *protocol.PlatformMessage {
	msg := BuildP2HDefaultMessage(host.ID, host.Name)
	msg.Plugin = &protocol.PluginMessage{
		Target: BuildInstanceDescriptor(target, host.ID),
		Notification: &protocol.NotificationMessage{
			Type:      eventType,
			Timestamp: time.Now().Unix(),
			Data:      payload,
			Error:     nil,
		},
	}
	return msg
}

func BuildCallPluginConfigChangeMessage(
	key string, originVal, newVal []string,
	host common_type.HostInfo,
	target common_type.IInstanceDescription,
) *protocol.PlatformMessage {
	msg := BuildP2HDefaultMessage(host.ID, host.Name)
	msg.Plugin = &protocol.PluginMessage{
		Target: BuildInstanceDescriptor(target, host.ID),
		Config: &protocol.ConfigurationMessage{
			ConfigChangeRequest: &protocol.ConfigurationMessage_ConfigurationChangeMessage{
				ConfigKey:   key,
				NewValue:    newVal,
				OriginValue: originVal,
			},
		},
	}
	return msg
}

func BuildLifecycleMessage(
	action protocol.ControlMessage_PluginActionType,
	host common_type.HostInfo,
	appID, instanceID, name, lang, langVer, appVer string,
	oldVersion *protocol.PluginDescriptor,
) *protocol.PlatformMessage {
	msg := BuildP2HDefaultMessage(host.ID, host.Name)
	msg.Control.LifeCycleRequest = &protocol.ControlMessage_PluginLifeCycleRequestMessage{
		Instance: &protocol.PluginInstanceDescriptor{
			Application: &protocol.PluginDescriptor{
				ApplicationID:      appID,
				Name:               name,
				Language:           lang,
				LanguageVersion:    VersionString2Pb(langVer),
				ApplicationVersion: VersionString2Pb(appVer),
				HostVersion:        VersionString2Pb(host.Version),
				MinSystemVersion:   VersionString2Pb(host.MinSystemVersion),
			},
			InstanceID: instanceID,
		},
		Action:     action,
		Reason:     "",
		OldVersion: oldVersion,
	}
	return msg
}

func BuildHostReportInitMessage(hostInfo *protocol.HostDescriptor) *protocol.PlatformMessage {
	msg := BuildH2PDefaultMessage(hostInfo.HostID, hostInfo.Name)
	msg.Control.HostReport = &protocol.ControlMessage_HostReportMessage{Host: hostInfo}
	return msg
}

func BuildHostBootReportInitMessage(bootInfo *protocol.HostBootDescriptor) *protocol.PlatformMessage {
	msg := BuildB2PDefaultMessage(bootInfo.BootID, bootInfo.Name)
	msg.Control.BootReport = &protocol.ControlMessage_HostBootReportMessage{Boot: bootInfo}
	return msg
}

func GetResourceInitMessage(source *protocol.PlatformMessage) *protocol.PlatformMessage {
	message := &protocol.PlatformMessage{
		Header: &protocol.RouterMessage{
			SeqNo:    source.Header.SeqNo,
			Source:   source.Header.Distinct,
			Distinct: source.Header.Source,
		},
		Resource: &protocol.ResourceMessage{
			Sender: source.Resource.Sender,
			Host:   source.Resource.Host,
		},
	}
	return message
}

func BuildP2HHeartbeatMessage(hostID, hostName string) *protocol.PlatformMessage {
	msg := BuildP2HDefaultMessage(hostID, hostName)
	msg.Control = &protocol.ControlMessage{Heartbeat: math.CreateCaptcha()}
	return msg
}

func BuildP2BHeartbeatMessage(bootID, bootName string) *protocol.PlatformMessage {
	msg := BuildP2BDefaultMessage(bootID, bootName)
	msg.Control = &protocol.ControlMessage{Heartbeat: math.CreateCaptcha()}
	return msg
}
