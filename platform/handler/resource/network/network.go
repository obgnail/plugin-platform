package network

import (
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/message_utils"
	"github.com/obgnail/plugin-platform/common/protocol"
)

type Network struct {
	source     *protocol.PlatformMessage
	distinct   *protocol.PlatformMessage
	opType     protocol.HttpResourceMessage_HttpResourceType
	instanceID string
}

func NewNetWork(sourceMessage, distinctMessage *protocol.PlatformMessage) *Network {
	network := &Network{
		source:   sourceMessage,
		distinct: distinctMessage,
	}
	network.instanceID = sourceMessage.GetResource().GetSender().GetInstanceID()
	network.opType = sourceMessage.GetResource().GetHttp().GetResourceType()
	return network
}

func (n *Network) Execute() {
	switch n.opType {
	case protocol.HttpResourceMessage_API:
		n.callMainSystem()
	case protocol.HttpResourceMessage_Outdoor:
		n.outdoor()
	}
}

func (n *Network) callMainSystem() {
	sender := n.source.GetResource().GetSender()
	HttpReqMsg := n.source.GetResource().GetHttp().GetResourceHttpRequest()

	msg := &protocol.HttpResponseMessage{}

	req := &CallMainSystemRequest{
		AppUUID:    sender.GetApplication().GetApplicationID(),
		InstanceID: n.instanceID,
		Method:     HttpReqMsg.GetMethod(),
		URL:        HttpReqMsg.GetUrl(),
		Header:     GetHeader(HttpReqMsg),
		PostData:   HttpReqMsg.GetBody(),
		Root:       HttpReqMsg.GetRoot(),
		Name:       sender.GetApplication().GetName(),
	}
	resp, err := CallMainSystemAPI(req)
	if err != nil {
		e := common_type.NewPluginError(common_type.CallMainSystemAPIFailure, err.Error(), common_type.CallMainSystemAPIFailureError.Error())
		msg.Error = message_utils.BuildErrorMessage(e)
	}

	if resp != nil {
		msg.Body = resp.Body
		msg.Headers = resp.BuildHeader()
		msg.StatusCode = int64(resp.StatusCode)
	}
	message_utils.BuildResourceNetworkMessage(n.distinct, msg)
}

func (n *Network) outdoor() {

}
