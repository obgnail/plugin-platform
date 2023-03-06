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
	"github.com/obgnail/plugin-platform/platform/conn/resource"
	resourceLog "github.com/obgnail/plugin-platform/platform/conn/resource/log"
	"time"
)

const (
	defaultTimeoutSec       = 30
	defaultHeartbeatSec     = 5
	defaultRetryIntervalSec = 1
)

var Timeout = time.Duration(config.Int("platform.timeout_sec", defaultTimeoutSec)) * time.Second
var RetryInterval = time.Duration(config.Int("platform.retry_interval_sec", defaultRetryIntervalSec)) * time.Second
var HeartbeatInterval = time.Duration(config.Int("platform.heartbeat_sec", defaultHeartbeatSec)) * time.Second

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
		h.onHostReport(msg)
	}

	if control.GetBootReport() != nil {
		h.onHostBootReport(msg)
	}
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

func (h *PlatformHandler) onHostBootReport(msg *protocol.PlatformMessage) {
	report := msg.GetControl().GetBootReport()
	if report == nil {
		return
	}

	log.Trace("【GET】HostBootReport. %+v", report)
	h.logMsg(report.GetLog())

	hostBootID := report.GetBoot().GetBootID()
	if hostBootID == "" {
		return
	}

	hostBoot := NewHostBoot(report, common_type.HostBootStatusNormal)
	h.hostBootPool.Add(hostBoot)
}

func (h *PlatformHandler) onHostReport(msg *protocol.PlatformMessage) {
	report := msg.GetControl().GetHostReport()
	if report == nil {
		return
	}

	log.Trace("【GET】HostReport. %+v", report)

	h.logMsg(report.GetLog())

	hostID := report.GetHost().GetHostID()
	if hostID == "" {
		return
	}
	_host := NewHost(report, common_type.HostStatusNormal)
	h.hostPool.Add(_host)
}

func (h *PlatformHandler) logMsg(logMsg []*protocol.LogMessage) {
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

func (h *PlatformHandler) lifeCycle(
	done chan common_type.PluginError,
	action protocol.ControlMessage_PluginActionType,
	appID, instanceID, name, lang, langVer, appVer string,
	oldVersion *protocol.PluginDescriptor,
) chan common_type.PluginError {
	h._getHostByInstanceID(instanceID, func(host common_type.IHost) {
		if host == nil {
			done <- common_type.NewPluginError(common_type.MsgTimeOut, "get host timeout")
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
		h.conn.SendAsync(msg, Timeout, func(input, result *protocol.PlatformMessage, err common_type.PluginError) {
			if err != nil {
				log.PEDetails(err)
				h.hostPool.DeleteByID(info.ID)
				log.Warn("delete host: %s", info.ID)
			} else {
				h.onHostReport(result) // 返回hostReport信息,这里需要及时更新
			}
			if done != nil {
				done <- err
			}
		})
	})
	return done
}

func (h *PlatformHandler) EnablePlugin(done chan common_type.PluginError, appID, instanceID, name, lang, langVer, appVer string) chan common_type.PluginError {
	return h.lifeCycle(done, protocol.ControlMessage_Enable, appID, instanceID, name, lang, langVer, appVer, nil)
}

func (h *PlatformHandler) DisablePlugin(done chan common_type.PluginError, appID, instanceID, name, lang, langVer, appVer string) chan common_type.PluginError {
	return h.lifeCycle(done, protocol.ControlMessage_Disable, appID, instanceID, name, lang, langVer, appVer, nil)
}

func (h *PlatformHandler) InstallPlugin(done chan common_type.PluginError, appID, instanceID, name, lang, langVer, appVer string) chan common_type.PluginError {
	return h.lifeCycle(done, protocol.ControlMessage_Install, appID, instanceID, name, lang, langVer, appVer, nil)
}

func (h *PlatformHandler) UnInstallPlugin(done chan common_type.PluginError, appID, instanceID, name, lang, langVer, appVer string) chan common_type.PluginError {
	return h.lifeCycle(done, protocol.ControlMessage_UnInstall, appID, instanceID, name, lang, langVer, appVer, nil)
}

func (h *PlatformHandler) UpgradePlugin(done chan common_type.PluginError, appID, instanceID, name, lang, langVer, appVer string, oldVersion *protocol.PluginDescriptor) chan common_type.PluginError {
	return h.lifeCycle(done, protocol.ControlMessage_Upgrade, appID, instanceID, name, lang, langVer, appVer, oldVersion)
}

func (h *PlatformHandler) CheckStatePlugin(done chan common_type.PluginError, appID, instanceID, name, lang, langVer, appVer string) chan common_type.PluginError {
	return h.lifeCycle(done, protocol.ControlMessage_CheckState, appID, instanceID, name, lang, langVer, appVer, nil)
}

func (h *PlatformHandler) CheckCompatibilityPlugin(done chan common_type.PluginError, appID, instanceID, name, lang, langVer, appVer string) chan common_type.PluginError {
	return h.lifeCycle(done, protocol.ControlMessage_CheckCompatibility, appID, instanceID, name, lang, langVer, appVer, nil)
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

func (h *PlatformHandler) KillPlugin(instanceID string) {
	host := h.GetHost(instanceID)
	if host == nil {
		return
	}
	info := host.GetInfo()

	log.Warn("kill plugin: host:%s. plugin:%s", info.ID, instanceID)

	msg := message.BuildP2HDefaultMessage(info.ID, info.Name)
	msg.Control.KillPlugin = &protocol.ControlMessage_KillPluginMessage{InstanceID: instanceID}
	h.conn.SendAsync(msg, Timeout, func(input, result *protocol.PlatformMessage, err common_type.PluginError) {
		if err == nil {
			log.Warn("kill plugin: %s", instanceID)
			return
		}
		log.PEDetails(err)
		h.hostPool.DeleteByID(info.ID)
		log.Warn("delete hostPool: %s", info.ID)
	})
}

func (h *PlatformHandler) GetAllHost() []common_type.IHost {
	return h.hostPool.GetAll()
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

func (h *PlatformHandler) _createHost() common_type.IHost {
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
	maxRetry := 5
	for {
		Host, ok := h.hostPool.Exist(result.Control.StartHost.Host.HostID)
		if ok {
			return Host
		}

		if count != maxRetry {
			count++
			time.Sleep(RetryInterval)
		} else {
			break
		}
	}
	return nil
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

func (h *PlatformHandler) GetHostBoot(hostBootID string) common_type.IHostBoot {
	hostboot, _ := h.hostBootPool.Exist(hostBootID)
	return hostboot
}

func (h *PlatformHandler) _getHostByInstanceID(instanceID string, callback func(host common_type.IHost)) {
	host := h.GetHost(instanceID)
	if host != nil {
		callback(host)
		return
	}

	go callback(h._createHost())
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
