package hub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/obgnail/plugin-platform/common/config"
	"github.com/obgnail/plugin-platform/common/errors"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	urlCallAbility = "/plugin/call_ability"
	urlCallEvent   = "/plugin/on_event"

	defaultTimeoutSec = 30
)

var (
	timeout      = time.Duration(config.Int("main_system.timeout_sec", defaultTimeoutSec)) * time.Second
	platformAddr = fmt.Sprintf("http://%s:%d",
		config.String("platform.host", "127.0.0.1"),
		config.Int("platform.http_port", 9005),
	)
)

func callPlugin(instanceID, abilityID, abilityType, abilityFuncKey string, arg []byte) ([]byte, error) {
	url := platformAddr + urlCallAbility

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(arg))
	if err != nil {
		return nil, errors.Trace(err)
	}
	req.Header.Set("instanceID", instanceID)
	req.Header.Set("abilityID", abilityID)
	req.Header.Set("abilityType", abilityType)
	req.Header.Set("abilityFunc", abilityFuncKey)

	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Trace(err)
	}
	defer resp.Body.Close()

	if resp == nil || resp.Body == nil {
		return nil, nil
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return respData, errors.Trace(err)
	}

	return respData, nil
}

type Resp struct {
	Code    int    `json:"code"`
	ErrCode string `json:"errcode"`
	Type    string `json:"type"`
}

func callEvent(instanceID, eventType string, payload []byte) error {
	url := platformAddr + urlCallEvent
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return errors.Trace(err)
	}
	req.Header.Set("instanceID", instanceID)
	req.Header.Set("eventType", eventType)
	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Trace(err)
	}
	defer resp.Body.Close()
	if resp == nil || resp.Body == nil {
		return nil
	}
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Trace(err)
	}
	fmt.Println("---", string(respData))
	rep := &Resp{}
	if err := json.Unmarshal(respData, rep); err != nil {
		return errors.Trace(err)
	}
	if rep.Code != http.StatusOK {
		return fmt.Errorf(rep.ErrCode)
	}
	return nil
}
