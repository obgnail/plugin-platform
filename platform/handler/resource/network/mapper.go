package network

import "github.com/obgnail/plugin-platform/common/protocol"

type PermissionMapper interface {
	Map(req *protocol.HttpRequestMessage, plugin *protocol.PluginInstanceDescriptor) *protocol.HttpRequestMessage
}

type DummyPermissionMapper struct{}

func (m *DummyPermissionMapper) Map(req *protocol.HttpRequestMessage, plugin *protocol.PluginInstanceDescriptor) *protocol.HttpRequestMessage {
	return req
}
