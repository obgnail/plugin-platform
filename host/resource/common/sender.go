package common

import (
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/connect"
	"github.com/obgnail/plugin-platform/common/protocol"
)

type Sender interface {
	Send(plugin common_type.IPlugin, msg *protocol.PlatformMessage) (resp *protocol.PlatformMessage, err common_type.PluginError)
	SendAsync(plugin common_type.IPlugin, msg *protocol.PlatformMessage, callback connect.CallBack)
	SendOnly(msg *protocol.PlatformMessage) (err common_type.PluginError)
}
