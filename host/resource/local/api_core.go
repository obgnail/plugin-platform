package local

import (
	"github.com/obgnail/plugin-platform/common/common_type"
)

var _ common_type.APICore = (*APICoreOp)(nil)

type APICoreOp struct {
	//msg    *protocol.HttpRequestMessage
	plugin common_type.IPlugin
}

func NewAPICore(plugin common_type.IPlugin) common_type.APICore {
	return &APICoreOp{plugin: plugin}
}

//
//func (n *APICoreOp) buildMessage(httpRequestMessage *protocol.HttpRequestMessage) *protocol.PlatformMessage {
//	msg := utils.GetInitMessage()
//	msg.Resource = &protocol.ResourceMessage{
//		Http: &protocol.HttpResourceMessage{
//			ResourceType:        protocol.HttpResourceMessage_API,
//			ResourceHttpRequest: httpRequestMessage,
//		},
//	}
//	return msg
//}
//
//func (n *APICoreOp) sendMsgToHost(platformMessage *protocol.PlatformMessage) (*protocol.PlatformMessage, common_type.PluginError) {
//	return SyncSendToHost(n.plugin, platformMessage)
//}

func (a *APICoreOp) Fetch(httpRequest *common_type.HttpRequest) *common_type.HttpResponse {
	//headers := make(map[string]*protocol.HeaderVal)
	//for k, val := range httpRequest.Headers {
	//	headerVal := &protocol.HeaderVal{Val: val}
	//	headers[k] = headerVal
	//}
	//httpRequestMessage := &protocol.HttpRequestMessage{
	//	Method:  httpRequest.Method,
	//	Url:     "http://127.0.0.1:9001" + httpRequest.Path,
	//	Body:    httpRequest.Body,
	//	Headers: headers,
	//	Root:    httpRequest.Root,
	//}
	//msg, err := a.sendMsgToHost(a.buildMessage(httpRequestMessage))
	//if err != nil {
	//	return &common_type.HttpResponse{
	//		Err:        err,
	//		StatusCode: 500,
	//	}
	//}
	//
	//body := msg.GetResource().GetHttp().GetResourceHttpResponse().GetBody()
	//headers = msg.GetResource().GetHttp().GetResourceHttpResponse().GetHeaders()
	//responseHeaders := make(map[string][]string)
	//for k, v := range headers {
	//	for _, val := range v.GetVal() {
	//		responseHeaders[k] = append(responseHeaders[k], val)
	//	}
	//}
	//
	//if reterr := msg.GetResource().GetHttp().GetResourceHttpResponse().GetError(); reterr != nil {
	//	return &common_type.HttpResponse{
	//		Err:        common_type.NewPluginError(int(reterr.Code), reterr.GetError(), reterr.GetMsg()),
	//		StatusCode: utils.Int64ToInt(msg.GetResource().GetHttp().GetResourceHttpResponse().GetStatusCode()),
	//		Headers:    responseHeaders,
	//		Body:       body,
	//	}
	//}
	//
	//return &common_type.HttpResponse{
	//	Headers:    responseHeaders,
	//	Body:       body,
	//	StatusCode: utils.Int64ToInt(msg.GetResource().GetHttp().GetResourceHttpResponse().GetStatusCode()),
	//}
	return nil
}
