package message_utils

import (
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/config"
	"github.com/obgnail/plugin-platform/common/connect"
	"github.com/obgnail/plugin-platform/common/math"
	"github.com/obgnail/plugin-platform/common/protocol"
)

var (
	hostID       = config.StringOrPanic("host.id")
	hostName     = config.StringOrPanic("host.name")
	platformID   = config.StringOrPanic("platform.id")
	platformName = config.StringOrPanic("platform.name")
)

func GetInitMessage() *protocol.PlatformMessage {
	message := &protocol.PlatformMessage{
		Header: &protocol.RouterMessage{
			SeqNo: math.CreateCaptcha(),
			Source: &protocol.RouterNode{
				Tags: make(map[string]string),
			},
			Distinct: &protocol.RouterNode{
				Tags: make(map[string]string),
			},
		},
		Control: &protocol.ControlMessage{},
	}
	return message
}

func BuildHostToPlatFormMessageWithHeader() *protocol.PlatformMessage {
	msg := GetInitMessage()
	msg.Header.Source = GetHostInfo()
	msg.Header.Distinct = GetPlatformInfo()
	msg.Header.SeqNo = math.CreateCaptcha()
	return msg
}

func GetHostInfo() *protocol.RouterNode {
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

func GetResourceInitMessage(sourceMessage *protocol.PlatformMessage) *protocol.PlatformMessage {
	message := &protocol.PlatformMessage{
		Header: &protocol.RouterMessage{
			Source: &protocol.RouterNode{
				Tags: make(map[string]string),
			},
			Distinct: &protocol.RouterNode{
				Tags: make(map[string]string),
			},
		},
		Resource: &protocol.ResourceMessage{},
	}
	swapMessageHeader(message, sourceMessage)
	message.Resource.Sender = sourceMessage.Resource.Sender
	message.Resource.Host = sourceMessage.Resource.Host
	return message
}

//消息头 来源/去向 调换
func swapMessageHeader(newMsg, comeMsg *protocol.PlatformMessage) {
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

func BuildResourceFileMessage(distinctMessage *protocol.PlatformMessage, resp *protocol.WorkspaceMessage_IOResponseMessage) {
	// workspaceMessage
	workspaceMessage := &protocol.WorkspaceMessage{}
	workspaceMessage.IOResponse = resp
	distinctMessage.Resource.Workspace = workspaceMessage
}

func BuildResourceDbMessage(distinctMessage *protocol.PlatformMessage, resp *protocol.DatabaseMessage_DatabaseResponseMessage) {
	// dataBaseMsgResp
	dataBaseMsg := &protocol.DatabaseMessage{}
	dataBaseMsg.DBResponse = resp
	distinctMessage.Resource.Database = dataBaseMsg
}

func BuildResourceNetworkMessage(distinctMessage *protocol.PlatformMessage, resp *protocol.HttpResponseMessage) {
	httpMessage := &protocol.HttpResourceMessage{}
	httpMessage.ResourceHttpResponse = resp
	distinctMessage.Resource.Http = httpMessage
}

func BuildResourceEventMessage(distinctMessage *protocol.PlatformMessage, resp *protocol.EventMessage) {
	distinctMessage.Resource.Event = resp
}

func BuildInstanceDescriptor(plugin common_type.IPlugin, hostID string) *protocol.PluginInstanceDescriptor {
	desc := plugin.GetPluginDescription().PluginDescription()
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
		InstanceID: plugin.GetPluginDescription().InstanceID(),
		HostID:     hostID,
	}
}

//func BuildReportMessage(newMsg, comeMsg *protocol.PlatformMessage, instanceList map[string]*protocol.PluginInstanceDescriptor, hostDesc *protocol.HostDescriptor) {
//	swapMessageHeader(newMsg, comeMsg)
//
//	reportMsg := &protocol.ControlMessage_HostReportMessage{
//		Host:         hostDesc,
//		InstanceList: instanceList,
//	}
//	control := &protocol.ControlMessage{}
//	control.Report = reportMsg
//	newMsg.Control = control
//}
//
//func BuildReportInitMessage(hostConfig *golangcommon.HostConfig, hostInfo *protocol.HostDescriptor) *protocol.PlatformMessage {
//	msg := GetInitMessage()
//	msg.Header.Source = GetHostInfo(hostConfig)
//	msg.Header.Distinct = GetPlatformInfo(hostConfig)
//	msg.Header.SeqNo = utils.CreateCaptcha()
//	msg.Control.Report = &protocol.ControlMessage_HostReportMessage{
//		Host: hostInfo,
//	}
//	return msg
//}
