package handler

import (
	"fmt"
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/config"
	"github.com/obgnail/plugin-platform/common/connect"
	"github.com/obgnail/plugin-platform/common/errors"
	"github.com/obgnail/plugin-platform/common/log"
	"github.com/obgnail/plugin-platform/common/protocol"
	"github.com/obgnail/plugin-platform/common/utils/math"
	"github.com/obgnail/plugin-platform/common/utils/message"
	"github.com/obgnail/plugin-platform/platform/conn/hub/event"
	"github.com/obgnail/plugin-platform/platform/conn/resource"
	resourceLog "github.com/obgnail/plugin-platform/platform/conn/resource/log"
	"time"
)

const (
	defaultTimeoutSec       = 30
	defaultHeartbeatSec     = 5
	defaultRetryIntervalSec = 1
	defaultMaxRetry         = 5
)

var Timeout = time.Duration(config.Int("platform.timeout_sec", defaultTimeoutSec)) * time.Second
var HeartbeatInterval = time.Duration(config.Int("platform.heartbeat_sec", defaultHeartbeatSec)) * time.Second
var RetryInterval = time.Duration(config.Int("platform.retry_interval_sec", defaultRetryIntervalSec)) * time.Second
var MaxRetry = config.Int("platform.max_retry", defaultMaxRetry)

var _ connect.ConnectionHandler = (*PlatformHandler)(nil)

type PlatformHandler struct {
	hostPool     *HostPool
	hostBootPool *HostBootPool
	conn         *connect.Connection // 负责和host、hostboot的通讯
}

func New(id, name, addr string) *PlatformHandler {
	h := &PlatformHandler{hostPool: NewHostPool(), hostBootPool: NewHostBootPool()}
	zmq := connect.NewZmq(id, name, addr, connect.SocketTypeRouter, connect.RolePlatform).SetPacker(&connect.ProtoPacker{})
	h.conn = connect.NewConnection(zmq, h)
	return h
}

func Default() *PlatformHandler {
	id := config.String("platform.id", "P0000001")
	name := config.String("platform.name", "platform")
	port := config.Int("platform.tcp_port", 9006)
	h := New(id, name, fmt.Sprintf("tcp://*:%d", port))

	log.Info("init Platform handler: ID:%s, Name:%s, Port:%d", id, name, port)

	return h
}

func (h *PlatformHandler) Heartbeat() {
	host := func() {
		for _, _host := range h.hostPool.GetAll() {
			info := _host.GetInfo()
			msg := message.BuildP2HHeartbeatMessage(info.ID, info.Name)
			h.conn.SendAsync(msg, Timeout, func(input, result *protocol.PlatformMessage, err common_type.PluginError) {
				if err == nil {
					return
				}
				log.PEDetails(err)
				h.hostPool.Delete(_host) // 过期失活
				log.Warn("delete host: %s", info.ID)
			})
		}
	}

	hostBoot := func() {
		for _, boot := range h.hostBootPool.GetAll() {
			info := boot.GetInfo()
			msg := message.BuildP2BHeartbeatMessage(info.ID, info.Name)
			h.conn.SendAsync(msg, Timeout, func(input, result *protocol.PlatformMessage, err common_type.PluginError) {
				if err == nil {
					return
				}
				log.PEDetails(err)
				h.hostBootPool.Delete(boot) // 过期失活
				log.Warn("delete hostBoot: %s", info.ID)
			})
		}
	}

	ticker := time.NewTicker(HeartbeatInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			host()
			hostBoot()
		}
	}
}

func (h *PlatformHandler) OnConnect() common_type.PluginError {
	log.Info("PlatformHandler OnConnect")
	return nil
}

func (h *PlatformHandler) OnDisconnect() common_type.PluginError {
	log.Info("PlatformHandler OnDisconnect")
	return nil
}

func (h *PlatformHandler) OnError(pluginError common_type.PluginError) {
	log.Error("%+v", pluginError)
}

func (h *PlatformHandler) OnMsg(endpoint *connect.EndpointInfo, msg *protocol.PlatformMessage, err common_type.PluginError) {
	if err != nil {
		log.PEDetails(err)
		return
	}

	h.OnControlMessage(endpoint, msg)
	h.OnResourceMessage(endpoint, msg)
}

// OnControlMessage 控制流
func (h *PlatformHandler) OnControlMessage(endpoint *connect.EndpointInfo, msg *protocol.PlatformMessage) {
	control := msg.GetControl()
	if control == nil {
		return
	}

	if control.GetHostReport() != nil {
		h.OnHostReport(msg)
	}

	if control.GetBootReport() != nil {
		h.OnHostBootReport(msg)
	}
}

func (h *PlatformHandler) OnHostReport(msg *protocol.PlatformMessage) {
	report := msg.GetControl().GetHostReport()
	if report == nil {
		return
	}

	log.Trace("【GET】HostReport. %+v", report)

	h.logMessage(report.GetLog())

	hostID := report.GetHost().GetHostID()
	if hostID == "" {
		return
	}
	_host := NewHost(report, common_type.HostStatusNormal)
	h.hostPool.Add(_host)
}

func (h *PlatformHandler) OnHostBootReport(msg *protocol.PlatformMessage) {
	report := msg.GetControl().GetBootReport()
	if report == nil {
		return
	}

	log.Trace("【GET】HostBootReport. %+v", report)
	h.logMessage(report.GetLog())

	hostBootID := report.GetBoot().GetBootID()
	if hostBootID == "" {
		return
	}

	hostBoot := NewHostBoot(report, common_type.HostBootStatusNormal)
	h.hostBootPool.Add(hostBoot)
}

// OnResourceMessage 资源
func (h *PlatformHandler) OnResourceMessage(endpoint *connect.EndpointInfo, msg *protocol.PlatformMessage) {
	req := msg.GetResource()
	if req == nil {
		return
	}

	log.Trace("【GET】GetResource. %+v", req)
	resp := resource.NewExecutor(msg).Execute()
	if resp != nil {
		log.Trace("【SND】GetResource. %+v", resp)
		if err := h.conn.SendOnly(resp); err != nil {
			log.PEDetails(err)
		}
	}
}

func (h *PlatformHandler) CallPluginConfigChanged(instanceID, configKey string, originValue, newValue []string) chan common_type.PluginError {
	errChan := make(chan common_type.PluginError, 1)

	msgBuilder := func(host common_type.HostInfo, plugin common_type.IInstanceDescription) *protocol.PlatformMessage {
		return message.BuildCallPluginConfigChangeMessage(configKey, originValue, newValue, host, plugin)
	}

	go func() {
		result, err := h.SendToAlivePlugin(instanceID, msgBuilder)
		if err != nil {
			errChan <- err
			return
		}

		if e := result.Plugin.Config.ConfigChangeResponse; e != nil {
			errChan <- common_type.NewPluginError(common_type.CallPluginHttpFailure, e.Msg)
		}
		errChan <- nil
	}()

	return errChan
}

func (h *PlatformHandler) CallPluginEvent(instanceID string, eventType string, payload []byte, force bool) chan common_type.PluginError {
	errChan := make(chan common_type.PluginError, 1)

	// force: 无论插件是否订阅了该topic, 强制通知
	if !force {
		ok := event.FilterEvent(instanceID, eventType, payload)
		if !ok {
			msg := fmt.Sprintf("instance %s does not subscribe this topic: %s", instanceID, eventType)
			errChan <- common_type.NewPluginError(common_type.NotifyEventFailure, msg)
			return errChan
		}
	}

	msgBuilder := func(host common_type.HostInfo, plugin common_type.IInstanceDescription) *protocol.PlatformMessage {
		return message.BuildCallPluginEventMessage(eventType, payload, host, plugin)
	}

	go func() {
		result, err := h.SendToAlivePlugin(instanceID, msgBuilder)
		if err != nil {
			errChan <- err
			return
		}
		if e := result.Plugin.Notification.Error; e != nil {
			errChan <- common_type.NewPluginError(common_type.NotifyEventFailure, e.Msg)
		}
		errChan <- nil
	}()
	return errChan
}

func (h *PlatformHandler) CallPluginFunction(instanceID string, abilityID, abilityType, abilityFuncKey string, args []byte) chan *common_type.AbilityResponse {
	respChan := make(chan *common_type.AbilityResponse, 1)

	msgBuilder := func(host common_type.HostInfo, plugin common_type.IInstanceDescription) *protocol.PlatformMessage {
		return message.BuildCallPluginFunctionMessage(abilityID, abilityType, abilityFuncKey, args, host, plugin)
	}

	go func() {
		result, err := h.SendToAlivePlugin(instanceID, msgBuilder)
		if err != nil {
			respChan <- &common_type.AbilityResponse{Err: err}
			return
		}

		respObj := result.Plugin.Ability.AbilityResponse
		resp := &common_type.AbilityResponse{Data: respObj.Data}
		if respErr := respObj.Error; respErr != nil {
			resp.Err = common_type.NewPluginError(common_type.CallAbilityFailure, respErr.Msg)
		}

		respChan <- resp
	}()

	return respChan
}

func (h *PlatformHandler) CallPluginHTTP(instanceID string, req *common_type.HttpRequest, internal bool, abilityFunc string) chan *common_type.HttpResponse {
	respChan := make(chan *common_type.HttpResponse, 1)

	msgBuilder := func(host common_type.HostInfo, plugin common_type.IInstanceDescription) *protocol.PlatformMessage {
		return message.BuildCallPluginHTTPMessage(req, internal, host, plugin, abilityFunc)
	}

	go func() {
		result, err := h.SendToAlivePlugin(instanceID, msgBuilder)
		if err != nil {
			respChan <- &common_type.HttpResponse{Err: err}
			return
		}

		respObj := result.Plugin.Http.Response

		var pe common_type.PluginError
		if respObj.Error != nil {
			pe = common_type.NewPluginError(common_type.CallPluginHttpFailure, respObj.Error.Msg)
		}

		headers := make(map[string][]string)
		for k, v := range respObj.Headers {
			for _, v1 := range v.Val {
				headers[k] = append(headers[k], v1)
			}
		}

		resp := &common_type.HttpResponse{
			Err:        pe,
			Headers:    headers,
			Body:       respObj.Body,
			StatusCode: int(respObj.StatusCode),
		}
		respChan <- resp
	}()

	return respChan
}

func (h *PlatformHandler) KillHost(hostID string) {
	host, ok := h.hostPool.Exist(hostID)
	if !ok {
		return
	}
	info := host.GetInfo()
	msg := message.BuildP2HDefaultMessage(info.ID, info.Name)
	msg.Control.Kill = &protocol.ControlMessage_KillPluginHostMessage{Kill: true}
	if err := h.conn.SendOnly(msg); err != nil {
		log.PEDetails(err)
	}
}

func (h *PlatformHandler) KillPlugin(instanceID string) chan common_type.PluginError {
	errChan := make(chan common_type.PluginError, 1)

	go func() {
		host := h.GetHost(instanceID)
		if host == nil {
			errChan <- common_type.NewPluginError(common_type.HostNotFoundErr, "get host error")
			return
		}

		info := host.GetInfo()
		log.Warn("kill plugin: host:%s. plugin:%s", info.ID, instanceID)

		msg := message.BuildP2HDefaultMessage(info.ID, info.Name)
		msg.Control.KillPlugin = &protocol.ControlMessage_KillPluginMessage{InstanceID: instanceID}

		_, err := h.conn.Send(msg, Timeout)
		if err != nil {
			log.PEDetails(err)
			h.hostPool.DeleteByID(info.ID)
			log.Warn("delete hostPool: %s", info.ID)
		} else {
			log.Warn("kill plugin: %s", instanceID)
		}
		errChan <- err
	}()
	return errChan
}

func (h *PlatformHandler) lifeCycle(
	action protocol.ControlMessage_PluginActionType,
	appID, instanceID, name, lang, langVer, appVer string,
	oldVersion *protocol.PluginDescriptor,
) chan common_type.PluginError {
	done := make(chan common_type.PluginError, 1)
	go func() {
		host := h.getHostByInstanceID(instanceID)
		if host == nil {
			done <- common_type.NewPluginError(common_type.MsgTimeOut, "get host error")
			return
		}

		info := host.GetInfo()
		msg := message.BuildLifecycleMessage(action, info, appID, instanceID, name, lang, langVer, appVer, oldVersion)
		result, err := h.conn.Send(msg, Timeout)
		if err != nil {
			log.PEDetails(err)
			h.hostPool.DeleteByID(info.ID)
			log.Warn("delete host: %s", info.ID)
		} else {
			h.OnHostReport(result) // 返回hostReport信息,这里需要及时更新
		}
		done <- err
	}()

	return done
}

// lifeCycleInSupported 调用已经运行的插件的生命周期
// (因为处于运行状态,所以在pool里面有相应的数据,不需要手动传递)
func (h *PlatformHandler) lifeCycleInSupported(action protocol.ControlMessage_PluginActionType,
	instanceID string, oldVersion *protocol.PluginDescriptor) chan common_type.PluginError {
	plugins := h.GetAllSupportPlugin()
	target, ok := plugins[instanceID]
	if !ok {
		done := make(chan common_type.PluginError, 1)
		done <- common_type.NewPluginError(common_type.GetInstanceFailure, "no such instance")
		return done
	}

	desc := target.PluginDescription()
	appID := desc.ApplicationID()
	name := desc.Name()
	lang := desc.Language()
	langVer := desc.LanguageVersion().VersionString()
	appVer := desc.ApplicationVersion().VersionString()
	return h.lifeCycle(action, appID, instanceID, name, lang, langVer, appVer, oldVersion)
}

func (h *PlatformHandler) InstallPlugin(appID, instanceID, name, lang, langVer, appVer string) chan common_type.PluginError {
	return h.lifeCycle(protocol.ControlMessage_Install, appID, instanceID, name, lang, langVer, appVer, nil)
}

func (h *PlatformHandler) UpgradePlugin(appID, instanceID, name, lang, langVer, appVer string, oldVersion *protocol.PluginDescriptor) chan common_type.PluginError {
	return h.lifeCycle(protocol.ControlMessage_Upgrade, appID, instanceID, name, lang, langVer, appVer, oldVersion)
}

func (h *PlatformHandler) EnablePlugin(instanceID string) chan common_type.PluginError {
	return h.lifeCycleInSupported(protocol.ControlMessage_Enable, instanceID, nil)
}

func (h *PlatformHandler) DisablePlugin(instanceID string) chan common_type.PluginError {
	return h.lifeCycleInSupported(protocol.ControlMessage_Disable, instanceID, nil)
}

func (h *PlatformHandler) UnInstallPlugin(instanceID string) chan common_type.PluginError {
	return h.lifeCycleInSupported(protocol.ControlMessage_UnInstall, instanceID, nil)
}

func (h *PlatformHandler) CheckStatePlugin(instanceID string) chan common_type.PluginError {
	return h.lifeCycleInSupported(protocol.ControlMessage_CheckState, instanceID, nil)
}

func (h *PlatformHandler) CheckCompatibilityPlugin(instanceID string) chan common_type.PluginError {
	return h.lifeCycleInSupported(protocol.ControlMessage_CheckCompatibility, instanceID, nil)
}

func (h *PlatformHandler) GetAllHost() []common_type.IHost {
	return h.hostPool.GetAll()
}

func (h *PlatformHandler) GetHost(instanceID string) common_type.IHost {
	var ret common_type.IHost
	h.hostPool.Range(func(hostID string, host common_type.IHost) bool {
		plugins := host.GetInfo().RunningPlugins
		for _instanceID := range plugins {
			if _instanceID == instanceID {
				ret = host
				return false
			}
		}
		return true
	})

	if ret == nil {
		h.hostPool.Range(func(hostID string, host common_type.IHost) bool {
			plugins := host.GetInfo().SupportPlugins
			for _instanceID := range plugins {
				if _instanceID == instanceID {
					ret = host
					return false
				}
			}
			return true
		})
	}

	return ret
}

func (h *PlatformHandler) CreateHost() common_type.IHost {
	boot := h.hostBootPool.One()
	if boot == nil {
		log.Error("has no host boot")
		return nil
	}

	info := boot.GetInfo()

	id := fmt.Sprintf("Host-%d", math.CreateCaptcha())
	name := id

	msg := message.BuildP2BDefaultMessage(info.ID, info.Name)
	msg.Control.StartHost = &protocol.ControlMessage_StartHostMessage{
		Host: &protocol.HostDescriptor{HostID: id, Name: name},
	}

	result, err := h.conn.Send(msg, Timeout)
	if err != nil {
		log.PEDetails(err)
		h.hostBootPool.Delete(boot)
		return nil
	}
	log.Info("start host: %+v", result.Control.StartHost)

	count := 0
	for {
		Host, ok := h.hostPool.Exist(result.Control.StartHost.Host.HostID)
		if ok {
			return Host
		}

		if count != MaxRetry {
			count++
			time.Sleep(RetryInterval)
		} else {
			break
		}
	}
	return nil
}

func (h *PlatformHandler) GetHostBoot(hostBootID string) common_type.IHostBoot {
	hostboot, _ := h.hostBootPool.Exist(hostBootID)
	return hostboot
}

func (h *PlatformHandler) GetAllHostBoot() []common_type.IHostBoot {
	return h.hostBootPool.GetAll()
}

func (h *PlatformHandler) GetAllAlivePlugin() map[string]common_type.IInstanceDescription {
	ret := make(map[string]common_type.IInstanceDescription)
	h.hostPool.Range(func(hostID string, host common_type.IHost) bool {
		for _, plugin := range host.GetInfo().RunningPlugins {
			ret[plugin.InstanceID()] = plugin
		}
		return true
	})
	return ret
}

func (h *PlatformHandler) GetAllSupportPlugin() map[string]common_type.IInstanceDescription {
	ret := make(map[string]common_type.IInstanceDescription)
	h.hostPool.Range(func(hostID string, host common_type.IHost) bool {
		for _, plugin := range host.GetInfo().SupportPlugins {
			ret[plugin.InstanceID()] = plugin
		}
		return true
	})
	return ret
}

type MessageBuilder func(host common_type.HostInfo, plugin common_type.IInstanceDescription) *protocol.PlatformMessage

func (h *PlatformHandler) SendToAlivePlugin(instanceID string, messageBuilder MessageBuilder) (*protocol.PlatformMessage, common_type.PluginError) {
	plugins := h.GetAllAlivePlugin()
	target, ok := plugins[instanceID]
	if !ok {
		return nil, common_type.NewPluginError(common_type.GetInstanceFailure, fmt.Sprintf("instanceNotFound: %s", instanceID))
	}

	host := h.getHostByInstanceID(instanceID)
	if host == nil {
		return nil, common_type.NewPluginError(common_type.MsgTimeOut, "get host timeout")
	}

	hostInfo := host.GetInfo()
	msg := messageBuilder(hostInfo, target)
	result, err := h.conn.Send(msg, Timeout)
	if err != nil {
		log.PEDetails(err)
	}
	return result, err
}

func (h *PlatformHandler) getHostByInstanceID(instanceID string) (host common_type.IHost) {
	if host = h.GetHost(instanceID); host != nil {
		return host
	}

	host = h.CreateHost()
	return host
}

func (h *PlatformHandler) logMessage(logMsg []*protocol.LogMessage) {
	if logMsg != nil {
		go func() {
			logger, err := resourceLog.NewLogger(config.StringOrPanic("platform.log_path"))
			if err != nil {
				log.ErrorDetails(errors.Trace(err))
				return
			}
			for _, l := range logMsg {
				logger.Log(l)
			}
		}()
	}
}

func (h *PlatformHandler) Run() common_type.PluginError {
	log.Info("PlatformHandler Run")

	err := h.conn.GetZmq().Connect()
	if err != nil {
		return err
	}
	go h.Heartbeat()
	return nil
}
