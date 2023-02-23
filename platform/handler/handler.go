package handler

import (
	"fmt"
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/config"
	"github.com/obgnail/plugin-platform/common/connect"
	"github.com/obgnail/plugin-platform/common/log"
	"github.com/obgnail/plugin-platform/common/protocol"
)

type PlatformHandler struct {
	*connect.BaseHandler
}

func New(id, name, addr string) *PlatformHandler {
	h := &PlatformHandler{}
	zmq := connect.NewZmq(id, name, addr, connect.SocketTypeRouter, connect.RolePlatform).SetPacker(&connect.ProtoPacker{})
	h.BaseHandler = connect.NewBaseHandler(zmq, h)
	return h
}

func Default() *PlatformHandler {
	id := config.String("platform.id", "R0000001")
	name := config.String("platform.name", "platform")
	port := config.Int("platform.tcp_port", 9006)
	h := New(id, name, fmt.Sprintf("tcp://*:%d", port))

	log.Info("init Platform handler: ID:%s, Name:%s, Port:%d", id, name, port)

	return h
}

func (ph *PlatformHandler) OnConnect() common_type.PluginError {
	log.Info("OnConnect")
	return nil
}

func (ph *PlatformHandler) OnDisconnect() common_type.PluginError {
	log.Info("OnDisconnect")
	return nil
}

func (ph *PlatformHandler) OnMsg(endpoint *connect.EndpointInfo, message *protocol.PlatformMessage, unmarshalError common_type.PluginError) {
	if unmarshalError != nil {
		log.ErrorDetails(unmarshalError)
		return
	}

	if message.GetResource() != nil {
		log.Info("message.GetResource() GetSeqNo: %d", message.GetHeader().GetSeqNo())
		//resourceOp := resource.NewResourceOp(message)
		//resourceOp.OnResource()
		//if resourceOp.DistinctMessage != nil {
		//	log.Logger.Info("message.GetResource() llll GetSeqNo:: %d", resourceOp.DistinctMessage.GetHeader().GetSeqNo())
		//	err = ph.SendMessage(resourceOp.DistinctMessage)
		//	if err != nil {
		//		ph.log.ErrorDetails(err)
		//		return
		//	}
		//}
	}
}

func (ph *PlatformHandler) OnError(pluginError common_type.PluginError) {
	log.ErrorDetails(pluginError)
}

func (ph *PlatformHandler) Run() common_type.PluginError {
	log.Info("Platform handler run")

	return ph.GetZmq().Connect()
}
