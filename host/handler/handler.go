package handler

import (
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/config"
	"github.com/obgnail/plugin-platform/common/connect"
	"github.com/obgnail/plugin-platform/common/log"
	"github.com/obgnail/plugin-platform/common/message_utils"
	"github.com/obgnail/plugin-platform/common/protocol"
	"os"
	"time"
)

// TODO
var defaultTimeout = 30 * time.Second

//var _ connect.MessageHandler = (*HostHandler)(nil)
var _ connect.FurtherHandler = (*HostHandler)(nil)

type HostHandler struct {
	descriptor *protocol.HostDescriptor
	*connect.BaseHandler
	isLocal bool
}

func New(id, name, addr, lang, hostVersion, minSysVersion, langVersion string, isLocal bool) *HostHandler {
	handler := &HostHandler{
		descriptor: &protocol.HostDescriptor{
			HostID:           id,
			Name:             name,
			Language:         lang,
			HostVersion:      message_utils.SplitVersion(hostVersion),
			MinSystemVersion: message_utils.SplitVersion(minSysVersion),
			LanguageVersion:  message_utils.SplitVersion(langVersion),
		},
		isLocal: isLocal,
	}

	zmq := connect.NewZmq(id, name, addr, connect.SocketTypeDealer, connect.RoleHost).SetPacker(&connect.ProtoPacker{})
	handler.BaseHandler = connect.NewBaseHandler(zmq, handler)
	return handler
}

func Default(isLocal bool) *HostHandler {
	id := config.StringOrPanic("host.id")
	name := config.StringOrPanic("host.name")
	addr := config.StringOrPanic("host.platform_address")
	lang := config.StringOrPanic("host.language")
	hostVersion := config.StringOrPanic("host.version")
	minSysVersion := config.StringOrPanic("host.min_system_version")
	langVersion := config.StringOrPanic("host.language_version")

	h := New(id, name, addr, lang, hostVersion, minSysVersion, langVersion, isLocal)
	return h
}

func (h *HostHandler) GetDescriptor() *protocol.HostDescriptor {
	return h.descriptor
}

func (h *HostHandler) OnConnect() common_type.PluginError {
	log.Info("OnConnect: %s", h.descriptor.Name)
	return nil
}

func (h *HostHandler) OnDisconnect() common_type.PluginError {
	log.Info("OnDisconnect: %s", h.descriptor.Name)
	return nil
}

func (h *HostHandler) OnError(err common_type.PluginError) {
	log.Warn("OnError: %s", h.descriptor.Name)
	if err.Code() != common_type.EndpointReceiveErr {
		os.Exit(1)
	}
	time.Sleep(time.Second * 9)
	if e := h.GetZmq().Connect(); e != nil {
		os.Exit(1)
	}
}

func (h *HostHandler) OnMsg(endpoint *connect.EndpointInfo, msg *protocol.PlatformMessage, err common_type.PluginError) {
	if err != nil {
		log.ErrorDetails(err)
		return
	}

	// 资源请求的应答
	if msg.Resource != nil {
		log.Info("%+v", msg)
	}
}

func (h *HostHandler) Send(sender common_type.IPlugin, msg *protocol.PlatformMessage) (*protocol.PlatformMessage, common_type.PluginError) {
	h.fillMsg(sender, msg)
	result, err := h.BaseHandler.Send(msg, defaultTimeout)
	return result, err
}

func (h *HostHandler) SendAsync(sender common_type.IPlugin, msg *protocol.PlatformMessage, callback connect.CallBack) {
	h.fillMsg(sender, msg)
	h.BaseHandler.SendAsync(msg, defaultTimeout, callback)
}

// fillMsg 添加路由信息
func (h *HostHandler) fillMsg(sender common_type.IPlugin, msg *protocol.PlatformMessage) {
	if msg.GetResource() != nil {
		msg.Resource.Sender = message_utils.BuildInstanceDescriptor(sender, h.descriptor.HostID)
		msg.Resource.Host = h.descriptor
		////补全日志信息
		//if msg.GetResource().GetLog() != nil {
		//	msg.Resource.Log.PluginInstanceDescriptor = pluginInstanceDescriptor
		//	msg.Resource.Log.HostDescriptor = hh.hostDescriptor
		//}
	}
}

func (h *HostHandler) Run() common_type.PluginError {
	return h.GetZmq().Connect()
}
