package handler

import (
	"fmt"
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/protocol"
	"github.com/obgnail/plugin-platform/common/utils/message"
	"reflect"
	"runtime"
)

const (
	defaultInternalOnCallFuncName     = "OnCall"
	defaultFuncSignatureDismatchError = "function signature dismatch"
)

type Caller interface {
	CallAbility(plugin common_type.IPlugin, req *protocol.StandardAbilityMessage_AbilityRequestMessage) (*protocol.StandardAbilityMessage_AbilityResponseMessage, common_type.PluginError)
	CallHTTP(plugin common_type.IPlugin, req *protocol.HttpRequestMessage) (*protocol.HttpResponseMessage, common_type.PluginError)
}

var _ Caller = (*pluginCaller)(nil)

type pluginCaller struct{}

func NewPluginCaller() *pluginCaller {
	return &pluginCaller{}
}

func (p *pluginCaller) CallAbility(plugin common_type.IPlugin, req *protocol.StandardAbilityMessage_AbilityRequestMessage) (
	*protocol.StandardAbilityMessage_AbilityResponseMessage, common_type.PluginError) {

	request := p.buildAbilityRequest(req)
	data, err := p.CallFunction(plugin, request)
	if err != nil {
		return nil, err
	}
	return p.convert2AbilityPb(data), nil
}

func (p *pluginCaller) buildAbilityRequest(req *protocol.StandardAbilityMessage_AbilityRequestMessage) *common_type.AbilityRequest {
	r := &common_type.AbilityRequest{
		ID:   req.Id,
		Type: req.Type,
		Func: req.FuncKey,
		Args: req.Args,
	}
	return r
}

func (p *pluginCaller) convert2AbilityPb(resp *common_type.AbilityResponse) *protocol.StandardAbilityMessage_AbilityResponseMessage {
	msg := &protocol.StandardAbilityMessage_AbilityResponseMessage{
		Data:  resp.Data,
		Error: message.BuildErrorMessage(resp.Err),
	}
	return msg
}

// CallFunction 规定abilityFunc的签名:
//    myFunc(req *common_type.AbilityRequest) *common_type.AbilityResponse
func (p *pluginCaller) CallFunction(plugin common_type.IPlugin, req *common_type.AbilityRequest) (
	*common_type.AbilityResponse, common_type.PluginError) {

	resp, e := CallFunction(plugin, req.Func, req)
	if e != nil {
		return nil, common_type.NewPluginError(common_type.CallAbilityFailure, e.Error())
	}

	if len(resp) != 1 {
		return nil, common_type.NewPluginError(common_type.CallAbilityFailure, defaultFuncSignatureDismatchError)
	}

	data := resp[0]
	if data == nil {
		return nil, nil
	}

	respObj, ok := data.(*common_type.AbilityResponse)
	if !ok {
		return nil, common_type.NewPluginError(common_type.CallAbilityFailure, defaultFuncSignatureDismatchError)
	}

	return respObj, nil
}

func (p *pluginCaller) CallHTTP(plugin common_type.IPlugin, pbReq *protocol.HttpRequestMessage) (
	pbResp *protocol.HttpResponseMessage, err common_type.PluginError) {

	req := p.getReqObj(pbReq)

	var resp *common_type.HttpResponse
	if pbReq.Internal == false {
		resp = plugin.OnExternalHttpRequest(req)
	} else {
		abilityFunc := p.getAbilityFunc(pbReq)
		resp, err = p.callHTTPFunction(plugin, abilityFunc, req)
	}

	pbResp = p.convert2HTTPPb(resp)
	return pbResp, nil
}

func (p *pluginCaller) getAbilityFunc(pbReq *protocol.HttpRequestMessage) string {
	abilityFunc := pbReq.AbilityFunc
	if abilityFunc == "" {
		abilityFunc = defaultInternalOnCallFuncName
	}
	return abilityFunc
}

func (p *pluginCaller) convert2HTTPPb(resp *common_type.HttpResponse) *protocol.HttpResponseMessage {
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

func (p *pluginCaller) getReqObj(request *protocol.HttpRequestMessage) *common_type.HttpRequest {
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

func (p *pluginCaller) callHTTPFunction(plugin common_type.IPlugin, funcName string, req *common_type.HttpRequest) (
	*common_type.HttpResponse, common_type.PluginError) {

	data, err := CallFunction(plugin, funcName, req)
	if err != nil {
		return nil, common_type.NewPluginError(common_type.CallPluginHttpFailure, err.Error())
	}
	if len(data) == 0 {
		return nil, nil
	}

	d := data[0]
	if d == nil {
		return nil, nil
	}

	respObj, ok := d.(*common_type.HttpResponse)
	if !ok {
		return nil, common_type.NewPluginError(common_type.CallPluginHttpFailure, defaultFuncSignatureDismatchError)
	}
	return respObj, nil
}

func CallFunction(any interface{}, funcName string, funcArgs ...interface{}) ([]interface{}, error) {
	errChan := make(chan error, 1)
	dataset := invoke(errChan, any, funcName, funcArgs...)
	if err := <-errChan; err != nil {
		return nil, err
	}

	var result []interface{}
	for _, data := range dataset {
		result = append(result, data.Interface())
	}
	return result, nil
}

func invoke(errChan chan error, any interface{}, name string, args ...interface{}) []reflect.Value {
	defer func() {
		if err := recover(); err != nil {
			buf := make([]byte, 10240)
			n := runtime.Stack(buf, false)
			stackInfo := fmt.Sprintf("%s", buf[:n])
			er := fmt.Errorf("plugin invoke panic:\nobject:\t%+v\nfunc:\t%s\nargs:\t%+v\nerr:\t%v\nstack:\t%s",
				any, name, args, err, stackInfo)
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
	dataset := ptr.Call(inputs)
	return dataset
}
