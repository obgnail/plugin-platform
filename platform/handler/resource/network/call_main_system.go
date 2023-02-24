package network

import (
	"bytes"
	"fmt"
	"github.com/obgnail/plugin-platform/common/config"
	"github.com/obgnail/plugin-platform/common/log"
	"github.com/obgnail/plugin-platform/common/protocol"
	"io/ioutil"
	"net/http"
)

type CallMainSystemRequest struct {
	AppUUID    string
	InstanceID string
	Method     string
	URL        string
	Header     map[string][]string
	PostData   []byte
	Root       bool //  插件可能会有不同的权限,root为false是普通权限,为true是根权限。
	Name       string
}

func GetHeader(HttpReqMsg *protocol.HttpRequestMessage) map[string][]string {
	result := make(map[string][]string)
	for k, val := range HttpReqMsg.GetHeaders() {
		result[k] = append(result[k], val.GetVal()...)
	}
	return result
}

type CallMainSystemResponse struct {
	InstanceID string
	StatusCode int
	Header     map[string][]string
	Body       []byte
}

func (r *CallMainSystemResponse) BuildHeader() map[string]*protocol.HeaderVal {
	var result = make(map[string]*protocol.HeaderVal)
	for header, newVal := range r.Header {
		if _, ok := result[header]; !ok {
			result[header] = &protocol.HeaderVal{Val: newVal}
		} else {
			result[header].Val = append(result[header].Val, newVal...)
		}
	}
	return result
}

// CallMainSystemAPI 调用主系统接口
func CallMainSystemAPI(request *CallMainSystemRequest) (*CallMainSystemResponse, error) {
	url := fmt.Sprintf("%s%s", config.StringOrPanic("main_system.addr"), request.URL)

	req, err := http.NewRequest(request.Method, url, bytes.NewReader(request.PostData))
	if err != nil {
		log.ErrorDetails(err)
		return nil, err
	}

	roleName := fmt.Sprintf("plugin:%s", request.Name)
	req.Header.Set("request-role", roleName) // 给插件的请求打上标记
	for k, v := range request.Header {
		for _, v1 := range v {
			req.Header.Set(k, v1)
		}
	}

	if request.Root {
		// TODO: 插件可能会有不同的权限,root为false是普通权限,为true是根权限。因为根权限和主系统密切相关，还没想好实现
	}

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		log.ErrorDetails(err)
		return nil, err
	}

	var headers = make(map[string][]string)
	for k, v := range resp.Header {
		headers[k] = append(headers[k], v...)
	}

	bodyData, err := ioutil.ReadAll(resp.Body)
	result := &CallMainSystemResponse{
		InstanceID: request.InstanceID,
		StatusCode: resp.StatusCode,
		Header:     headers,
		Body:       bodyData,
	}
	return result, nil
}
