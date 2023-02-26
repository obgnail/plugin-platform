package local

import (
	"github.com/obgnail/plugin-platform/common/common_type"
)

var _ common_type.Network = (*NetworkOp)(nil)

type NetworkOp struct {
	//msg    *protocol.HttpRequestMessage
	plugin common_type.IPlugin
}

func NewNetwork(plugin common_type.IPlugin) common_type.Network {
	return &NetworkOp{plugin: plugin}
}

//
//func (n *NetworkOp) buildMessage(resourceHttpRequest *protocol.HttpRequestMessage) *protocol.PlatformMessage {
//	msg := utils.GetInitMessage()
//	msg.Resource = &protocol.ResourceMessage{
//		Http: &protocol.HttpResourceMessage{
//			ResourceType:        protocol.HttpResourceMessage_Outdoor,
//			ResourceHttpRequest: resourceHttpRequest,
//		},
//	}
//	return msg
//}
//
//func (n *NetworkOp) sendMsgToHost(platformMessage *protocol.PlatformMessage) (*protocol.PlatformMessage, common_type.PluginError) {
//	return SyncSendToHost(n.Plugin, platformMessage)
//}

func (n *NetworkOp) Fetch(httpRequest *common_type.HttpRequest) *common_type.HttpResponse {
	//headers := make(map[string]*protocol.HeaderVal)
	//for k, val := range httpRequest.Headers {
	//	headerVal := &protocol.HeaderVal{Val: val}
	//	headers[k] = headerVal
	//}
	//httpRequestMessage := &protocol.HttpRequestMessage{
	//	Method:  httpRequest.Method,
	//	Url:     httpRequest.Url + httpRequest.Path,
	//	Body:    httpRequest.Body,
	//	Headers: headers,
	//	Root:    httpRequest.Root,
	//}
	//msg, err := n.sendMsgToHost(n.buildMessage(httpRequestMessage))
	//if err != nil {
	//	return &common_type.HttpResponse{
	//		StatusCode: 500,
	//		Err:        err,
	//	}
	//}
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
	//		Headers:    responseHeaders,
	//		StatusCode: utils.Int64ToInt(msg.GetResource().GetHttp().GetResourceHttpResponse().GetStatusCode()),
	//		Err:        common_type.NewPluginError(int(reterr.Code), reterr.GetError(), reterr.GetMsg()),
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

//type AsyncNetwork struct {
//	callBackHandler common_type.NetworkCallBack
//	timeoutHandler  common_type.AsyncInvokeTimeoutCallback
//}
//
//func (a *AsyncNetwork) callBack(param, ret *protocol.PlatformMessage, asyncObject interface{}, err common_type.PluginError) {
//	httpResp := ret.GetResource().GetHttp().GetResourceHttpResponse()
//	resp := &common_type.HttpResponse{
//		Headers:    make(map[string][]string),
//		Body:       httpResp.Body,
//		StatusCode: int(httpResp.StatusCode),
//	}
//	for k, v := range httpResp.Headers {
//		resp.Headers[k] = v.Val
//	}
//	a.callBackHandler(resp, err, asyncObject)
//}
//
//func (a *AsyncNetwork) timeOutHandler(param *protocol.PlatformMessage, asyncObject interface{}, err common_type.PluginError) {
//	a.timeoutHandler(err, asyncObject)
//}

func (n *NetworkOp) AsyncFetch(httpRequest *common_type.HttpRequest, callback common_type.NetworkCallBack) {

	//headers := make(map[string]*protocol.HeaderVal)
	//for k, val := range httpRequest.Headers {
	//	headerVal := &protocol.HeaderVal{Val: val}
	//	headers[k] = headerVal
	//}
	//httpRequestMessage := &protocol.HttpRequestMessage{
	//	Method:  httpRequest.Method,
	//	Url:     httpRequest.Url + httpRequest.Path,
	//	Body:    httpRequest.Body,
	//	Headers: headers,
	//	Root:    httpRequest.Root,
	//}
	//asyncNetwork := new(AsyncNetwork)
	//asyncNetwork.callBackHandler = callback
	//asyncNetwork.timeoutHandler = timeoutHandler
	//AsyncSendToHost(n.Plugin, n.buildMessage(httpRequestMessage), object, asyncNetwork.callBack, asyncNetwork.timeOutHandler)
	return
}
