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
	"github.com/obgnail/plugin-platform/platform/handler/resource"
	resourceLog "github.com/obgnail/plugin-platform/platform/handler/resource/log"
	"time"
)

const (
	defaultTimeoutSec   = 30
	defaultHeartbeatSec = 5
)

var Timeout = time.Duration(config.Int("platform.timeout_sec", defaultTimeoutSec)) * time.Second

type Handler struct {
	hostPool     *HostPool
	hostBootPool *HostBootPool
	*connect.BaseHandler
}

func New(id, name, addr string) *Handler {
	h := &Handler{hostPool: NewHostPool(), hostBootPool: NewHostBootPool()}
	zmq := connect.NewZmq(id, name, addr, connect.SocketTypeRouter, connect.RolePlatform).SetPacker(&connect.ProtoPacker{})
	h.BaseHandler = connect.NewBaseHandler(zmq, h)
	return h
}

func Default() *Handler {
	id := config.String("platform.id", "P0000001")
	name := config.String("platform.name", "platform")
	port := config.Int("platform.tcp_port", 9006)
	h := New(id, name, fmt.Sprintf("tcp://*:%d", port))

	log.Info("init Platform handler: ID:%s, Name:%s, Port:%d", id, name, port)

	return h
}

func (h *Handler) OnConnect() common_type.PluginError {
	log.Info("PlatformHandler OnConnect")
	return nil
}

func (h *Handler) OnDisconnect() common_type.PluginError {
	log.Info("PlatformHandler OnDisconnect")
	return nil
}

func (h *Handler) OnError(pluginError common_type.PluginError) {
	log.ErrorDetails(pluginError)
}

func (h *Handler) logSendSyncError(send *protocol.PlatformMessage, err common_type.PluginError) {
	if err.Code() != common_type.MsgTimeOut {
		log.Error("Timeout: %+v", send)
	} else {
		log.Error("%+v", send)
	}
}

func (h *Handler) Heartbeat() {
	host := func() {
		for _, _host := range h.hostPool.GetAll() {
			info := _host.GetInfo()
			msg := message.BuildP2HHeartbeatMessage(info.ID, info.Name)
			h.SendAsync(msg, Timeout, func(input, result *protocol.PlatformMessage, err common_type.PluginError) {
				if err == nil {
					return
				}
				h.logSendSyncError(msg, err)
				h.hostPool.Delete(_host) // 过期失活
				log.Warn("delete hostPool: %s", info.ID)
			})
		}
	}

	hostBoot := func() {
		for _, boot := range h.hostBootPool.GetAll() {
			info := boot.GetInfo()
			msg := message.BuildP2BHeartbeatMessage(info.ID, info.Name)
			h.SendAsync(msg, Timeout, func(input, result *protocol.PlatformMessage, err common_type.PluginError) {
				if err == nil {
					return
				}
				h.logSendSyncError(msg, err)
				h.hostBootPool.Delete(boot) // 过期失活
				log.Warn("delete hostBootPool: %s", info.ID)
			})
		}
	}

	sec := config.Int("platform.heartbeat_sec", defaultHeartbeatSec)
	ticker := time.NewTicker(time.Second * time.Duration(sec))
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			host()
			hostBoot()
		}
	}
}

func (h *Handler) OnMsg(endpoint *connect.EndpointInfo, msg *protocol.PlatformMessage, err common_type.PluginError) {
	if err != nil {
		log.ErrorDetails(err)
		return
	}

	h.OnControlMessage(endpoint, msg)
	h.OnResourceMessage(endpoint, msg)
}

// OnControlMessage 控制流
func (h *Handler) OnControlMessage(endpoint *connect.EndpointInfo, msg *protocol.PlatformMessage) {
	control := msg.GetControl()
	if control == nil {
		return
	}

	if control.GetHostReport() != nil {
		h.onHostReport(msg)
	}

	if control.GetBootReport() != nil {
		h.onHostBootReport(msg)
	}
}

// OnResourceMessage 资源
func (h *Handler) OnResourceMessage(endpoint *connect.EndpointInfo, msg *protocol.PlatformMessage) {
	if msg.GetResource() == nil {
		return
	}

	log.Trace("【GET】message.GetResource() GetSeqNo: %d", msg.GetHeader().GetSeqNo())
	resp := resource.NewExecutor(msg).Execute()
	if resp != nil {
		log.Trace("【SEND】message.GetResource() GetSeqNo: %d", resp.GetHeader().GetSeqNo())
		if err := h.SendOnly(resp); err != nil {
			log.ErrorDetails(errors.Trace(err))
		}
	}
}

func (h *Handler) onHostBootReport(msg *protocol.PlatformMessage) {
	report := msg.GetControl().GetBootReport()
	if report == nil {
		return
	}

	log.Trace("【GET】message.HostBootReport. GetSeqNo: %d", msg.GetHeader().GetSeqNo())
	h.logMsg(report.GetLog())

	hostBootID := report.GetBoot().GetBootID()
	if hostBootID == "" {
		return
	}

	hostBoot := NewHostBoot(report, common_type.HostBootStatusNormal)
	h.hostBootPool.Add(hostBoot)
}

func (h *Handler) onHostReport(msg *protocol.PlatformMessage) {
	report := msg.GetControl().GetHostReport()
	if report == nil {
		return
	}

	log.Trace("【GET】message.HostReport. %+v", report)

	h.logMsg(report.GetLog())

	hostID := report.GetHost().GetHostID()
	if hostID == "" {
		return
	}
	_host := NewHost(report, common_type.HostStatusNormal)
	h.hostPool.Add(_host)
}

func (h *Handler) logMsg(logMsg []*protocol.LogMessage) {
	if logMsg != nil {
		go func() {
			logger, err := resourceLog.NewLogger(config.StringOrPanic("platform.log_path"))
			if err != nil {
				log.ErrorDetails(err)
				return
			}
			for _, l := range logMsg {
				logger.Log(l)
			}
		}()
	}
}

func (h *Handler) lifeCycleReq(
	action protocol.ControlMessage_PluginActionType,
	appID, instanceID, name, lang, langVer, appVer string,
	oldVersion *protocol.PluginDescriptor,
) {

	h._getHostByInstanceID(instanceID, func(host common_type.IHost) {
		if host == nil {
			log.Error("get host timeout!")
			return
		}

		info := host.GetInfo()

		msg := message.BuildP2HDefaultMessage(info.ID, info.Name)
		msg.Control.LifeCycleRequest = &protocol.ControlMessage_PluginLifeCycleRequestMessage{
			Instance: &protocol.PluginInstanceDescriptor{
				Application: &protocol.PluginDescriptor{
					ApplicationID:      appID,
					Name:               name,
					Language:           lang,
					LanguageVersion:    message.VersionString2Pb(langVer),
					ApplicationVersion: message.VersionString2Pb(appVer),
					HostVersion:        message.VersionString2Pb(info.Version),
					MinSystemVersion:   message.VersionString2Pb(info.MinSystemVersion),
				},
				InstanceID: instanceID,
			},
			Action:     action,
			Reason:     "",
			OldVersion: oldVersion,
		}
		h.SendAsync(msg, Timeout, func(input, result *protocol.PlatformMessage, err common_type.PluginError) {
			if err == nil {
				return
			}
			h.logSendSyncError(input, err)
			h.hostPool.DeleteByID(info.ID)
			log.Warn("delete hostPool: %s", info.ID)
		})
	})
}

func (h *Handler) EnablePlugin(appID, instanceID, name, lang, langVer, appVer string) {
	h.lifeCycleReq(protocol.ControlMessage_Enable, appID, instanceID, name, lang, langVer, appVer, nil)
}

func (h *Handler) DisablePlugin(appID, instanceID, name, lang, langVer, appVer string) {
	h.lifeCycleReq(protocol.ControlMessage_Disable, appID, instanceID, name, lang, langVer, appVer, nil)
}

func (h *Handler) StartPlugin(appID, instanceID, name, lang, langVer, appVer string) {
	h.lifeCycleReq(protocol.ControlMessage_Start, appID, instanceID, name, lang, langVer, appVer, nil)
}

func (h *Handler) StopPlugin(appID, instanceID, name, lang, langVer, appVer string) {
	h.lifeCycleReq(protocol.ControlMessage_Stop, appID, instanceID, name, lang, langVer, appVer, nil)
}

func (h *Handler) InstallPlugin(appID, instanceID, name, lang, langVer, appVer string) {
	h.lifeCycleReq(protocol.ControlMessage_Install, appID, instanceID, name, lang, langVer, appVer, nil)
}

func (h *Handler) UnInstallPlugin(appID, instanceID, name, lang, langVer, appVer string) {
	h.lifeCycleReq(protocol.ControlMessage_UnInstall, appID, instanceID, name, lang, langVer, appVer, nil)
}

func (h *Handler) UpgradePlugin(appID, instanceID, name, lang, langVer, appVer string, oldVersion *protocol.PluginDescriptor) {
	h.lifeCycleReq(protocol.ControlMessage_Upgrade, appID, instanceID, name, lang, langVer, appVer, oldVersion)
}

func (h *Handler) CheckStatePlugin(appID, instanceID, name, lang, langVer, appVer string) {
	h.lifeCycleReq(protocol.ControlMessage_CheckState, appID, instanceID, name, lang, langVer, appVer, nil)
}

func (h *Handler) CheckCompatibilityPlugin(appID, instanceID, name, lang, langVer, appVer string) {
	h.lifeCycleReq(protocol.ControlMessage_CheckCompatibility, appID, instanceID, name, lang, langVer, appVer, nil)
}

func (h *Handler) KillHost(hostID string) {
	host, ok := h.hostPool.Exist(hostID)
	if !ok {
		return
	}
	info := host.GetInfo()
	msg := message.BuildP2HDefaultMessage(info.ID, info.Name)
	msg.Control.Kill = &protocol.ControlMessage_KillPluginHostMessage{Kill: true}
	if err := h.SendOnly(msg); err != nil {
		log.ErrorDetails(err)
	}
}

func (h *Handler) KillPlugin(instanceID string) {
	host := h._findHostByPluginInstanceID(instanceID)
	if host == nil {
		return
	}
	info := host.GetInfo()

	log.Warn("kill plugin: host:%s. plugin:%s", info.ID, instanceID)

	msg := message.BuildP2HDefaultMessage(info.ID, info.Name)
	msg.Control.KillPlugin = &protocol.ControlMessage_KillPluginMessage{InstanceID: instanceID}
	h.SendAsync(msg, Timeout, func(input, result *protocol.PlatformMessage, err common_type.PluginError) {
		if err == nil {
			log.Warn("kill plugin: %s", instanceID)
			return
		}
		h.logSendSyncError(input, err)
		h.hostPool.DeleteByID(info.ID)
		log.Warn("delete hostPool: %s", info.ID)
	})
}

func (h *Handler) GetAllHost() []common_type.IHost {
	return h.hostPool.GetAll()
}

func (h *Handler) GetAllHostBoot() []common_type.IHostBoot {
	return h.hostBootPool.GetAll()
}

func (h *Handler) GetAllAlivePlugin() map[string]common_type.IInstanceDescription {
	ret := make(map[string]common_type.IInstanceDescription)
	h.hostPool.Range(func(hostID string, host common_type.IHost) bool {
		for _, plugin := range host.GetInfo().RunningPlugins {
			ret[plugin.InstanceID()] = plugin
		}
		return true
	})
	return ret
}

func (h *Handler) GetAllSupportPlugin() map[string]common_type.IInstanceDescription {
	ret := make(map[string]common_type.IInstanceDescription)
	h.hostPool.Range(func(hostID string, host common_type.IHost) bool {
		for _, plugin := range host.GetInfo().SupportPlugins {
			ret[plugin.InstanceID()] = plugin
		}
		return true
	})
	return ret
}

func (h *Handler) _createHost() common_type.IHost {
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

	result, err := h.Send(msg, Timeout)
	if err != nil {
		h.logSendSyncError(msg, err)
		h.hostBootPool.Delete(boot)
		return nil
	}
	log.Info("start host: %+v", result.Control.StartHost)

	count := 0
	maxRetry := 5
	for {
		Host, ok := h.hostPool.Exist(result.Control.StartHost.Host.HostID)
		if ok {
			return Host
		}

		if count != maxRetry {
			count++
			time.Sleep(time.Second)
		} else {
			break
		}
	}
	return nil
}

func (h *Handler) _findHostByPluginInstanceID(instanceID string) common_type.IHost {
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

func (h *Handler) _getHostByInstanceID(instanceID string, callback func(host common_type.IHost)) {
	host := h._findHostByPluginInstanceID(instanceID)
	if host != nil {
		callback(host)
		return
	}

	go func() {
		_host := h._createHost()
		callback(_host)
	}()
}

func (h *Handler) Run() common_type.PluginError {
	log.Info("PlatformHandler Run")

	err := h.GetZmq().Connect()
	if err != nil {
		return err
	}
	go h.Heartbeat()
	return nil
}
