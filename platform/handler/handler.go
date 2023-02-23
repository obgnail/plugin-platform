package handler

import (
	"fmt"
	"github.com/BangWork/ones-platform-api/app/handler/message_handlers/resource"
	common "github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/log"
	"github.com/obgnail/plugin-platform/common/protocol"
	"github.com/obgnail/plugin-platform/platform/config"
	"github.com/obgnail/plugin-platform/utils/connect"
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
	id := config.StringOrPanic("id")
	name := config.StringOrPanic("name")
	addr := fmt.Sprintf("tcp://*:%d", config.IntOrPanic("tcp_port"))
	h := New(id, name, addr)
	return h
}

func (ph *PlatformHandler) OnConnect() common.PluginError {
	log.Info("OnConnect")
	return nil
}

func (ph *PlatformHandler) OnDisconnect() common.PluginError {
	log.Info("OnDisconnect")
	return nil
}

func (ph *PlatformHandler) OnMsg(endpoint *connect.EndpointInfo, message *protocol.PlatformMessage, unmarshalError common.PluginError) {
	if unmarshalError != nil {
		log.ErrorDetails(unmarshalError)
		return
	}

	if message.GetResource() != nil {
		log.Info("message.GetResource() GetSeqNo: %d", message.GetHeader().GetSeqNo())
		resourceOp := resource.NewResourceOp(message)
		resourceOp.OnResource()
		if resourceOp.DistinctMessage != nil {
			log.Logger.Info("message.GetResource() llll GetSeqNo:: %d", resourceOp.DistinctMessage.GetHeader().GetSeqNo())
			err = ph.SendMessage(resourceOp.DistinctMessage)
			if err != nil {
				ph.log.ErrorDetails(err)
				return
			}
		}
	}
}

func (ph *PlatformHandler) OnError(pluginError common.PluginError) {
	log.ErrorDetails(pluginError)
}

func (ph *PlatformHandler) Run() common.PluginError {
	return ph.GetZmq().Connect()
}
