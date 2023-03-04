package connect

import (
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/protocol"
)

type Role string

const (
	RolePlatform = "platform"
	RoleHost     = "host"
	RoleHostBoot = "host_boot"
)

type SocketType string

const (
	SocketTypeRouter = "router"
	SocketTypeDealer = "dealer"
)

// MessagePacker 规定在ZmqEndpoint传输的数据中必须包含元数据(即:发送端和接收端的相关信息),使用 MessagePacker 将其组装或分离
type MessagePacker interface {
	// Unpack 从rawData中剥离出发送端和接收端,返回processedData
	Unpack(rawData []byte) (source, target *EndpointInfo, processedData []byte, err common_type.PluginError)
	// Pack 给出发送端和接收端,为rawData加上发送端和接收端信息,生成新的发送内容
	Pack(source, target *EndpointInfo, rawData []byte) (processedData []byte, err common_type.PluginError)
}

type MessageHandler interface {
	OnConnect() common_type.PluginError
	OnDisconnect() common_type.PluginError
	OnMessage(endpoint *EndpointInfo, content []byte)
	OnError(pluginError common_type.PluginError) // EndpointReceiveErr、EndpointIdentifyErr、EndpointSendErr only
}

type ConnectionHandler interface {
	OnConnect() common_type.PluginError
	OnDisconnect() common_type.PluginError
	OnMsg(endpoint *EndpointInfo, content *protocol.PlatformMessage, unmarshalError common_type.PluginError)
	OnError(pluginError common_type.PluginError)
}
