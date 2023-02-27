package release

import (
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/message_utils"
	"github.com/obgnail/plugin-platform/common/protocol"
	"github.com/obgnail/plugin-platform/host/resource/common"
)

var _ common_type.APICore = (*APICore)(nil)
var _ common_type.Network = (*Outdoor)(nil)

type CommonNetwork struct {
	Type   protocol.HttpResourceMessage_HttpResourceType
	msg    *protocol.HttpRequestMessage
	plugin common_type.IPlugin
	sender common.Sender
}

func newNetworkCommon(plugin common_type.IPlugin, sender common.Sender, Type protocol.HttpResourceMessage_HttpResourceType) *CommonNetwork {
	return &CommonNetwork{plugin: plugin, sender: sender, Type: Type}
}

func (c *CommonNetwork) sendMsgToHost(platformMessage *protocol.PlatformMessage) (*protocol.PlatformMessage, common_type.PluginError) {
	return c.sender.Send(c.plugin, platformMessage)
}

func (c *CommonNetwork) sendToHostAsync(platformMessage *protocol.PlatformMessage, callback common_type.NetworkCallBack) {
	cb := &networkCallbackWrapper{Func: callback}
	c.sender.SendAsync(c.plugin, platformMessage, cb.callBack)
}

func (c *CommonNetwork) buildMessage(httpRequestMessage *protocol.HttpRequestMessage) *protocol.PlatformMessage {
	msg := message_utils.BuildHostToPlatFormMessageWithHeader()
	msg.Resource = &protocol.ResourceMessage{
		Http: &protocol.HttpResourceMessage{
			ResourceType:        c.Type,
			ResourceHttpRequest: httpRequestMessage,
		},
	}
	return msg
}

func (c *CommonNetwork) fetch(platformMessage *protocol.PlatformMessage) *common_type.HttpResponse {
	msg, err := c.sendMsgToHost(platformMessage)
	if err != nil {
		return &common_type.HttpResponse{Err: err, StatusCode: 500}
	}

	resp := msg.GetResource().GetHttp().GetResourceHttpResponse()
	body := resp.GetBody()
	headers := resp.GetHeaders()
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

type networkCallbackWrapper struct {
	Func common_type.NetworkCallBack
}

func (w *networkCallbackWrapper) callBack(input, result *protocol.PlatformMessage, err common_type.PluginError) {
	httpResp := result.GetResource().GetHttp().GetResourceHttpResponse()
	resp := &common_type.HttpResponse{
		Headers:    make(map[string][]string),
		Body:       httpResp.Body,
		StatusCode: int(httpResp.StatusCode),
	}
	for k, v := range httpResp.Headers {
		resp.Headers[k] = v.Val
	}
	w.Func(resp, err)
}

type APICore struct {
	common *CommonNetwork
}

func NewAPICore(plugin common_type.IPlugin, sender common.Sender) *APICore {
	return &APICore{common: newNetworkCommon(plugin, sender, protocol.HttpResourceMessage_API)}
}

func (a *APICore) buildMsg(httpRequest *common_type.HttpRequest) *protocol.PlatformMessage {
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
	msg := a.common.buildMessage(httpRequestMessage)
	return msg
}

func (a *APICore) Fetch(httpRequest *common_type.HttpRequest) *common_type.HttpResponse {
	msg := a.buildMsg(httpRequest)
	return a.common.fetch(msg)
}

type Outdoor struct {
	common *CommonNetwork
}

func NewOutdoor(plugin common_type.IPlugin, sender common.Sender) *Outdoor {
	return &Outdoor{common: newNetworkCommon(plugin, sender, protocol.HttpResourceMessage_Outdoor)}
}

func (o *Outdoor) buildMsg(httpRequest *common_type.HttpRequest) *protocol.PlatformMessage {
	headers := make(map[string]*protocol.HeaderVal)
	for key, val := range httpRequest.Headers {
		headers[key] = &protocol.HeaderVal{Val: val}
	}

	httpRequestMessage := &protocol.HttpRequestMessage{
		Method:  httpRequest.Method,
		Url:     httpRequest.Url + httpRequest.Path,
		Headers: headers,
		Body:    httpRequest.Body,
		Root:    httpRequest.Root,
	}
	msg := o.common.buildMessage(httpRequestMessage)
	return msg
}

func (o *Outdoor) Fetch(httpRequest *common_type.HttpRequest) *common_type.HttpResponse {
	msg := o.buildMsg(httpRequest)
	return o.common.fetch(msg)
}

func (o *Outdoor) AsyncFetch(httpRequest *common_type.HttpRequest, callback common_type.NetworkCallBack) {
	msg := o.buildMsg(httpRequest)
	o.common.sendToHostAsync(msg, callback)
}
