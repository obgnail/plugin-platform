package handler

import (
	"github.com/golang/protobuf/proto"
	common "github.com/obgnail/plugin-platform/common_type"
	"github.com/obgnail/plugin-platform/utils/connect"
	"github.com/obgnail/plugin-platform/utils/protocol"
)

type platformHandler struct {
	endpoint *connect.ZmqEndpoint
}

func (ph *platformHandler) OnConnect() common.PluginError { return nil }

func (ph *platformHandler) OnDisconnect() common.PluginError { return nil }

func (ph *platformHandler) OnMessage(endpoint *connect.EndpointInfo, content []byte) {
	message := &protocol.PlatformMessage{}
	err := proto.Unmarshal(content, message)
	if err != nil {
		return
	}

	if message.GetResource() != nil {
		//log.Logger.Info("message.GetResource() ssss GetSeqNo: %d", message.GetHeader().GetSeqNo())
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

func (ph *platformHandler) OnError(err error) {

}
