package network

import (
	"bytes"
	"fmt"
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/config"
	"github.com/obgnail/plugin-platform/common/errors"
	"github.com/obgnail/plugin-platform/common/log"
	"github.com/obgnail/plugin-platform/common/message_utils"
	"github.com/obgnail/plugin-platform/common/protocol"
	"io/ioutil"
	"net/http"
	"strings"
)

type Network struct {
	source   *protocol.PlatformMessage
	distinct *protocol.PlatformMessage
	mapper   PermissionMapper
	outdoor  bool
}

func NewNetWork(sourceMessage, distinctMessage *protocol.PlatformMessage) *Network {
	outdoor := false
	switch sourceMessage.GetResource().GetHttp().GetResourceType() {
	case protocol.HttpResourceMessage_API:
		outdoor = false
	case protocol.HttpResourceMessage_Outdoor:
		outdoor = true
	}
	network := &Network{
		source:   sourceMessage,
		distinct: distinctMessage,
		outdoor:  outdoor,
		mapper:   &DummyPermissionMapper{},
	}
	return network
}

const (
	RolePluginHeader = "request-role"
)

func (n *Network) Execute() {
	sender := n.source.GetResource().GetSender()
	HttpReq := n.source.GetResource().GetHttp().GetResourceHttpRequest()

	// 插件可能会有不同的权限,root为false是普通权限,为true是根权限。
	// 因为获得根权限和主系统密切相关,所以将其委托给mapper实现
	// mapper通过修改http请求参数,使之获得根权限
	if root := HttpReq.GetRoot(); root && n.mapper != nil {
		HttpReq = n.mapper.Map(HttpReq, sender)
	}

	method := strings.ToUpper(HttpReq.GetMethod())
	url := HttpReq.GetUrl()
	data := HttpReq.GetBody()

	appUUID := sender.GetApplication().GetApplicationID()
	instanceID := sender.GetInstanceID()

	reqHeaders := make(map[string][]string)
	for k, val := range HttpReq.GetHeaders() {
		reqHeaders[k] = append(reqHeaders[k], val.GetVal()...)
	}

	if !n.outdoor {
		url = fmt.Sprintf("%s%s", config.StringOrPanic("main_system.addr"), url)

		// 给插件的请求打上标记
		val := fmt.Sprintf("plugin:%s.%s", appUUID, instanceID)
		reqHeaders[RolePluginHeader] = append(reqHeaders[RolePluginHeader], val)
	}

	respObj, err := Request(method, url, data, reqHeaders)
	if err != nil {
		n.buildMsg(0, nil, nil, err)
		return
	}
	defer respObj.Body.Close()

	code, respHeaders, body, e := GetFromResp(respObj)
	n.buildMsg(code, body, respHeaders, e)
}

func (n *Network) buildMsg(status int64, body []byte, headers map[string]*protocol.HeaderVal, err error) {
	msg := &protocol.HttpResponseMessage{
		Body:       body,
		Headers:    headers,
		StatusCode: status,
	}
	if err != nil {
		log.ErrorDetails(err)
		e := common_type.NewPluginError(common_type.CallMainSystemAPIFailure, err.Error(),
			common_type.CallMainSystemAPIFailureError.Error())
		msg.Error = message_utils.BuildErrorMessage(e)
	}
	message_utils.BuildResourceNetworkMessage(n.distinct, msg)
}

func Request(method, url string, data []byte, headers map[string][]string) (*http.Response, error) {
	reqObj, err := http.NewRequest(method, url, bytes.NewReader(data))
	if err != nil {
		return nil, errors.Trace(err)
	}
	for k, v := range headers {
		for _, v1 := range v {
			reqObj.Header.Set(k, v1)
		}
	}

	respObj, err := new(http.Client).Do(reqObj)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return respObj, nil
}

func GetFromResp(respObj *http.Response) (code int64, headers map[string]*protocol.HeaderVal, body []byte, err error) {
	code = int64(respObj.StatusCode)
	headers = make(map[string]*protocol.HeaderVal)
	for k, v := range respObj.Header {
		if _, ok := headers[k]; !ok {
			headers[k] = &protocol.HeaderVal{Val: v}
		} else {
			headers[k].Val = append(headers[k].Val, v...)
		}
	}
	body, err = ioutil.ReadAll(respObj.Body)
	if err != nil {
		return code, headers, body, errors.Trace(err)
	}
	return code, headers, body, nil
}
