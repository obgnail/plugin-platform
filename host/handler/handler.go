package handler

import (
	"github.com/golang/protobuf/proto"
	common "github.com/obgnail/plugin-platform/common_type"
	"github.com/obgnail/plugin-platform/utils/connect"
	"github.com/obgnail/plugin-platform/utils/protocol"
	"os"
	"time"
)

var _ connect.MessageHandler = (*HostHandler)(nil)

type HostHandler struct {
	descriptor *protocol.HostDescriptor

	endpoint *connect.ZmqEndpoint

	isLocal bool
}

func New(id, name, addr, lang, hostVersion, minSysVersion, langVersion string, isLocal bool) *HostHandler {
	handler := &HostHandler{
		descriptor: &protocol.HostDescriptor{
			HostID:           id,
			Name:             name,
			Language:         lang,
			HostVersion:      protocol.SplitVersion(hostVersion),
			MinSystemVersion: protocol.SplitVersion(minSysVersion),
			LanguageVersion:  protocol.SplitVersion(langVersion),
		},
		isLocal: isLocal,
	}

	handler.endpoint = connect.New(id, name, addr, connect.SocketTypeDealer, connect.RoleHost).
		SetPacker(&connect.ProtoPacker{}).
		SetHandler(handler)

	return handler
}

func (h *HostHandler) OnConnect() common.PluginError    { return nil }
func (h *HostHandler) OnDisconnect() common.PluginError { return nil }
func (h *HostHandler) OnError(err common.PluginError) {
	if err.Code() != common.EndpointReceiveErr {
		os.Exit(1)
	}
	time.Sleep(time.Second * 9)
	if e := h.endpoint.Connect(); e != nil {
		os.Exit(1)
	}
}

func (h *HostHandler) OnMessage(endpoint *connect.EndpointInfo, content []byte) {
	message := &protocol.PlatformMessage{}
	if err := proto.Unmarshal(content, message); err != nil {
		// TODO 处理异常
		return
	}

	// 资源请求的应答
	if message.Resource != nil {

	}
}

func (h *HostHandler) Run() common.PluginError {
	return h.endpoint.Connect()
}
