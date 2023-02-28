package message_utils

import (
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/config"
	"github.com/obgnail/plugin-platform/common/connect"
	"github.com/obgnail/plugin-platform/common/math"
	"github.com/obgnail/plugin-platform/common/protocol"
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

// host -> platform
func BuildP2HDefaultMessage(hostID, hostName string) *protocol.PlatformMessage {
	p := GetPlatformInfo()
	h := GetHostInfo(hostID, hostName)
	msg := GetInitMessage(p, h)
	return msg
}

// platform -> host
func BuildH2PDefaultMessage(hostID, hostName string) *protocol.PlatformMessage {
	p := GetPlatformInfo()
	h := GetHostInfo(hostID, hostName)
	msg := GetInitMessage(h, p)
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

func BuildHostReportInitMessage(hostInfo *protocol.HostDescriptor) *protocol.PlatformMessage {
	msg := BuildH2PDefaultMessage(hostInfo.HostID, hostInfo.Name)
	msg.Control.Report = &protocol.ControlMessage_HostReportMessage{Host: hostInfo}
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

func BuildHostReportMessage(source *protocol.PlatformMessage, instances map[string]*protocol.PluginInstanceDescriptor,
	hostDesc *protocol.HostDescriptor,
) *protocol.PlatformMessage {
	resp := &protocol.PlatformMessage{
		Header: &protocol.RouterMessage{
			SeqNo:    source.Header.SeqNo,
			Source:   source.Header.Distinct,
			Distinct: source.Header.Source,
		},
		Control: &protocol.ControlMessage{
			Report: &protocol.ControlMessage_HostReportMessage{
				Host:         hostDesc,
				InstanceList: instances,
			},
		},
	}
	return resp
}

func BuildLifeCycleResponseMessage(source *protocol.PlatformMessage) *protocol.PlatformMessage {
	resp := &protocol.PlatformMessage{
		Header: &protocol.RouterMessage{
			SeqNo:    source.Header.SeqNo,
			Source:   source.Header.Distinct,
			Distinct: source.Header.Source,
		},
		Control: &protocol.ControlMessage{
			LifeCycleResponse: &protocol.ControlMessage_PluginLifeCycleResponseMessage{
				Host:     source.Control.LifeCycleRequest.Host,
				Instance: source.Control.LifeCycleRequest.Instance,
				Result:   true, // 这个值后面可能会被修改
				Error:    nil,  // 这个值后面可能会被修改
			},
		},
	}
	return resp
}

func BuildPlatform2HostHeartbeatMessage(hostID, hostName string) *protocol.PlatformMessage {
	msg := BuildP2HDefaultMessage(hostID, hostName)
	msg.Control = &protocol.ControlMessage{Heartbeat: math.CreateCaptcha()}
	return msg
}
