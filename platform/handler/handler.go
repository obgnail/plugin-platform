package handler

import (
	"fmt"
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/config"
	"github.com/obgnail/plugin-platform/common/connect"
	"github.com/obgnail/plugin-platform/common/log"
	"github.com/obgnail/plugin-platform/common/protocol"
	"github.com/obgnail/plugin-platform/platform/handler/resource"
	"time"
)

//TODO
var defaultTimeout = 30 * time.Second

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
	log.Info("PlatformHandler OnConnect")
	return nil
}

func (ph *PlatformHandler) OnDisconnect() common_type.PluginError {
	log.Info("PlatformHandler OnDisconnect")
	return nil
}

func (ph *PlatformHandler) OnMsg(endpoint *connect.EndpointInfo, msg *protocol.PlatformMessage, err common_type.PluginError) {
	if err != nil {
		log.ErrorDetails(err)
		return
	}

	if msg.GetResource() != nil {
		log.Info("【GET】message.GetResource() GetSeqNo: %d", msg.GetHeader().GetSeqNo())
		resp := resource.NewExecutor(msg).Execute()
		if resp != nil {
			log.Info("【SEND】message.GetResource() GetSeqNo: %d", resp.GetHeader().GetSeqNo())
			if err := ph.SendOnly(resp); err != nil {
				log.ErrorDetails(err)
			}
		}
	}
}

func (ph *PlatformHandler) OnError(pluginError common_type.PluginError) {
	log.ErrorDetails(pluginError)
}

func (ph *PlatformHandler) Run() common_type.PluginError {
	log.Info("PlatformHandler Run")

	return ph.GetZmq().Connect()
}
