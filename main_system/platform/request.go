package platform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/obgnail/plugin-platform/common/errors"
	"io/ioutil"
	"net/http"
)

func Get(api string, reqObj interface{}) (resData []byte, err error) {
	return Request(http.MethodGet, api, reqObj)
}

func Post(api string, reqObj interface{}) (resData []byte, err error) {
	return Request(http.MethodPost, api, reqObj)
}

func Request(method string, url string, reqObj interface{}) (resData []byte, err error) {
	var requestBody []byte
	if reqObj != nil {
		requestBody, err = json.Marshal(reqObj)
		if err != nil {
			err = errors.Trace(err)
			return
		}
	}

	client := new(http.Client)
	var resp *http.Response
	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, errors.Trace(err)
	}
	if reqObj != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err = client.Do(req)
	if err != nil {
		return nil, errors.Trace(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		resData, _ = ioutil.ReadAll(resp.Body)
		return nil, errors.Trace(fmt.Errorf("call plugin api error: url=\"%s\", status=%d, body=\"%s\"", url, resp.StatusCode, resData))
	}

	if resp != nil && resp.Body != nil {
		resData, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return resData, errors.Trace(err)
		}
	}
	return resData, nil
}
