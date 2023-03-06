package handler

import (
	"fmt"
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/log"
	"github.com/obgnail/plugin-platform/common/protocol"
	"github.com/obgnail/plugin-platform/common/utils/message"
	"reflect"
	"runtime"
)

const (
	defaultOnCallFuncName = "OnCall"
)

type PluginCaller struct{}

func NewCaller() *PluginCaller {
	return &PluginCaller{}
}

func (p *PluginCaller) CallPlugin(plugin common_type.IPlugin, pbReq *protocol.HttpRequestMessage) (
	*protocol.HttpResponseMessage, error) {
	abilityFunc := pbReq.AbilityFunc
	if abilityFunc == "" {
		abilityFunc = defaultOnCallFuncName
	}

	log.Trace("on call function: %s", abilityFunc)

	req := p.getReqObj(pbReq)
	resp, err := p.httpRequest(plugin, abilityFunc, req)
	if err != nil {
		return nil, err
	}
	pbResp := p.getPbResp(resp)
	return pbResp, nil
}

func (p *PluginCaller) getPbResp(resp *common_type.HttpResponse) *protocol.HttpResponseMessage {
	if resp == nil {
		return nil
	}
	respMsg := &protocol.HttpResponseMessage{
		StatusCode: int64(resp.StatusCode),
		Error:      message.BuildErrorMessage(resp.Err),
		Body:       resp.Body,
	}
	for k, v := range resp.Headers {
		respMsg.Headers[k] = &protocol.HeaderVal{Val: v}
	}
	return respMsg
}

func (p *PluginCaller) getReqObj(request *protocol.HttpRequestMessage) *common_type.HttpRequest {
	req := &common_type.HttpRequest{
		Url:     request.Url,
		Method:  request.Method,
		Body:    request.Body,
		Headers: make(map[string][]string),
	}
	for k, v := range request.Headers {
		for _, v1 := range v.Val {
			req.Headers[k] = append(req.Headers[k], v1)
		}
	}
	return req
}

func (p *PluginCaller) httpRequest(plugin common_type.IPlugin, funcName string, req *common_type.HttpRequest) (
	*common_type.HttpResponse, error) {
	errChan := make(chan error, 1)
	data := invoke(errChan, plugin, funcName, req)
	if err := <-errChan; err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, nil
	}
	resp := data[0]
	respObj := resp.Interface().(*common_type.HttpResponse)
	return respObj, nil
}

func invoke(errChan chan error, any interface{}, name string, args ...interface{}) []reflect.Value {
	defer func() {
		if err := recover(); err != nil {
			buf := make([]byte, 10240)
			n := runtime.Stack(buf, false)
			stackInfo := fmt.Sprintf("%s", buf[:n])
			er := fmt.Errorf("plugin invoke panic: %v, %s", err, stackInfo)
			errChan <- er
		} else {
			errChan <- nil
		}
	}()

	inputs := make([]reflect.Value, len(args))
	for i := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	ptr := reflect.ValueOf(any).MethodByName(name)
	data := ptr.Call(inputs)
	return data
}
