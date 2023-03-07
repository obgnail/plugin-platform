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

// ResourceFactor 资源工厂,负责获取到local版本或release版本的资源
type ResourceFactor interface {
	GetResource(plugin common_type.IPlugin, sender Sender) common_type.IResources
}
