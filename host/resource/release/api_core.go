package release

import (
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/message_utils"
	"github.com/obgnail/plugin-platform/common/protocol"
	"github.com/obgnail/plugin-platform/host/handler"
)

type APICore struct {
	msg     *protocol.HttpRequestMessage
	plugin  common_type.IPlugin
	handler *handler.HostHandler
}

func NewAPICore(plugin common_type.IPlugin, handler *handler.HostHandler) common_type.APICore {
	return &APICore{plugin: plugin, handler: handler}
}

func (a *APICore) sendMsgToHost(platformMessage *protocol.PlatformMessage) (*protocol.PlatformMessage, common_type.PluginError) {
	return a.handler.Send(a.plugin, platformMessage)
}

func (a *APICore) buildMessage(httpRequestMessage *protocol.HttpRequestMessage) *protocol.PlatformMessage {
	msg := message_utils.GetHostToPlatFormMessage()
	msg.Resource = &protocol.ResourceMessage{
		Sender: message_utils.BuildInstanceDescriptor(a.plugin, a.handler.GetDescriptor().HostID),
		Http: &protocol.HttpResourceMessage{
			ResourceType:        protocol.HttpResourceMessage_API,
			ResourceHttpRequest: httpRequestMessage,
		},
	}
	return msg
}

func (a *APICore) Fetch(httpRequest *common_type.HttpRequest) *common_type.HttpResponse {
	headers := make(map[string]*protocol.HeaderVal)
	for key, val := range httpRequest.Headers {
		headers[key] = &protocol.HeaderVal{Val: val}
	}

	httpRequestMessage := &protocol.HttpRequestMessage{
		Method:  httpRequest.Method,
		Url:     httpRequest.Path,
		Headers: headers,
		Body:    httpRequest.Body,
		Root:    httpRequest.Root,
	}
	msg, err := a.sendMsgToHost(a.buildMessage(httpRequestMessage))
	if err != nil {
		return &common_type.HttpResponse{Err: err, StatusCode: 500}
	}

	resp := msg.GetResource().GetHttp().GetResourceHttpResponse()
	body := resp.GetBody()
	headers = resp.GetHeaders()
	code := resp.GetStatusCode()
	retErr := resp.GetError()
	responseHeaders := make(map[string][]string)
	for k, v := range headers {
		for _, val := range v.GetVal() {
			responseHeaders[k] = append(responseHeaders[k], val)
		}
	}
	if retErr != nil {
		return &common_type.HttpResponse{
			Err:        common_type.NewPluginError(int(retErr.Code), retErr.GetError(), retErr.GetMsg()),
			StatusCode: int(code),
			Headers:    responseHeaders,
			Body:       body,
		}

	}
	result := &common_type.HttpResponse{Headers: responseHeaders, Body: body, StatusCode: int(code)}
	return result
}
