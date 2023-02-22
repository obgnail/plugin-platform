package local

import (
	common "github.com/obgnail/plugin-platform/common_type"
)

var _ common.Network = (*NetworkOp)(nil)

type NetworkOp struct {
	//msg    *protocol.HttpRequestMessage
	plugin common.IPlugin
}

func NewNetwork(plugin common.IPlugin) common.Network {
	return &NetworkOp{plugin: plugin}
}

//
//func (n *NetworkOp) BuildMessage(resourceHttpRequest *protocol.HttpRequestMessage) *protocol.PlatformMessage {
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
//func (n *NetworkOp) SendMsgToHost(platformMessage *protocol.PlatformMessage) (*protocol.PlatformMessage, common.PluginError) {
//	return SyncSendToHost(n.plugin, platformMessage)
//}

func (n *NetworkOp) Fetch(httpRequest *common.HttpRequest) *common.HttpResponse {
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
	//msg, err := n.SendMsgToHost(n.BuildMessage(httpRequestMessage))
	//if err != nil {
	//	return &common.HttpResponse{
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
	//	return &common.HttpResponse{
	//		Headers:    responseHeaders,
	//		StatusCode: utils.Int64ToInt(msg.GetResource().GetHttp().GetResourceHttpResponse().GetStatusCode()),
	//		Err:        common.NewPluginError(int(reterr.Code), reterr.GetError(), reterr.GetMsg()),
	//		Body:       body,
	//	}
	//}
	//
	//return &common.HttpResponse{
	//	Headers:    responseHeaders,
	//	Body:       body,
	//	StatusCode: utils.Int64ToInt(msg.GetResource().GetHttp().GetResourceHttpResponse().GetStatusCode()),
	//}
	return nil
}

//type AsyncNetwork struct {
//	callBackHandler common.NetworkCallBack
//	timeoutHandler  common.AsyncInvokeTimeoutCallback
//}
//
//func (a *AsyncNetwork) callBack(param, ret *protocol.PlatformMessage, asyncObject interface{}, err common.PluginError) {
//	httpResp := ret.GetResource().GetHttp().GetResourceHttpResponse()
//	resp := &common.HttpResponse{
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
//func (a *AsyncNetwork) timeOutHandler(param *protocol.PlatformMessage, asyncObject interface{}, err common.PluginError) {
//	a.timeoutHandler(err, asyncObject)
//}

func (n *NetworkOp) AsyncFetch(httpRequest *common.HttpRequest,
	object interface{}, callback common.NetworkCallBack, timeoutHandler common.AsyncInvokeTimeoutCallback) {

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
	//AsyncSendToHost(n.plugin, n.BuildMessage(httpRequestMessage), object, asyncNetwork.callBack, asyncNetwork.timeOutHandler)
	return
}