package message

import (
	"github.com/obgnail/plugin-platform/common/config"
	"github.com/obgnail/plugin-platform/common/connect"
	"github.com/obgnail/plugin-platform/common/protocol"
	"github.com/obgnail/plugin-platform/common/seq"
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
			SeqNo: seq.CreateCaptcha(),
			Source: &protocol.RouterNode{
				Tags: make(map[string]string),
			},
			Distinct: &protocol.RouterNode{
				Tags: make(map[string]string),
			},
		},
		//Control: &protocol.ControlMessage{},
	}
	return message
}

func GetHostToPlatFormMessage() *protocol.PlatformMessage {
	msg := GetInitMessage()
	msg.Header.Source = GetHostInfo()
	msg.Header.Distinct = GetPlatformInfo()
	msg.Header.SeqNo = seq.CreateCaptcha()
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

//
////消息头 来源/去向 调换
//func swapMessageHeader(newMsg, comeMsg *protocol.PlatformMessage) {
//	newMsg.Header.Source.ID = comeMsg.Header.Distinct.ID
//	newMsg.Header.Source.Tags = comeMsg.Header.Distinct.Tags
//	newMsg.Header.Distinct.ID = comeMsg.Header.Source.ID
//	newMsg.Header.Distinct.Tags = comeMsg.Header.Source.Tags
//	newMsg.Header.SeqNo = comeMsg.Header.SeqNo
//}
//
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
