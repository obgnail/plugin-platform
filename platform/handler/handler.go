package handler

import (
	"fmt"
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/config"
	"github.com/obgnail/plugin-platform/common/connect"
	"github.com/obgnail/plugin-platform/common/errors"
	"github.com/obgnail/plugin-platform/common/log"
	"github.com/obgnail/plugin-platform/common/message_utils"
	"github.com/obgnail/plugin-platform/common/protocol"
	"github.com/obgnail/plugin-platform/platform/handler/host"
	"github.com/obgnail/plugin-platform/platform/handler/resource"
	log2 "github.com/obgnail/plugin-platform/platform/handler/resource/log"
	"time"
)

const (
	defaultTimeoutSec   = 30
	defaultHeartbeatSec = 5
)

var Timeout = time.Duration(config.Int("platform.timeout_sec", defaultTimeoutSec)) * time.Second

type PlatformHandler struct {
	pool *host.Pool
	*connect.BaseHandler
}

func New(id, name, addr string) *PlatformHandler {
	h := &PlatformHandler{
		pool: host.NewPool(),
	}
	zmq := connect.NewZmq(id, name, addr, connect.SocketTypeRouter, connect.RolePlatform).SetPacker(&connect.ProtoPacker{})
	h.BaseHandler = connect.NewBaseHandler(zmq, h)
	return h
}

func Default() *PlatformHandler {
	id := config.String("platform.id", "R0000001")
	name := config.String("platform.name", "platform")
	port := config.Int("platform.tcp_port", 9006)
	h := New(id, name, fmt.Sprintf("tcp://*:%d", port))

	log.Info("init Platform handler: ID:%s, Name:%s, Port:%d", id, name, port)

	return h
}

func (ph *PlatformHandler) OnConnect() common_type.PluginError {
	log.Info("PlatformHandler OnConnect")
	return nil
}

func (ph *PlatformHandler) OnDisconnect() common_type.PluginError {
	log.Info("PlatformHandler OnDisconnect")
	return nil
}

func (ph *PlatformHandler) OnError(pluginError common_type.PluginError) {
	log.ErrorDetails(pluginError)
}

func (ph *PlatformHandler) Heartbeat2Host() {
	heartbeat := func() {
		hosts := ph.pool.GetAll()
		for _, _host := range hosts {
			info := _host.GetInfo()
			msg := message_utils.BuildPlatform2HostHeartbeatMessage(info.ID, info.Name)
			ph.SendAsync(msg, Timeout, func(input, result *protocol.PlatformMessage, err common_type.PluginError) {
				if err == nil {
					return
				}
				log.Error("host %s deleted.", info.ID)
				if err.Code() != common_type.MsgTimeOut {
					log.Error("Timeout: %+v", msg)
				} else {
					log.ErrorDetails(err)
				}
				ph.pool.Delete(_host) // 过期失活
			})
		}
	}

	sec := config.Int("platform.timeout_sec", defaultHeartbeatSec)
	ticker := time.NewTicker(time.Second * time.Duration(sec))
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			heartbeat()
		}
	}
}

func (ph *PlatformHandler) OnMsg(endpoint *connect.EndpointInfo, msg *protocol.PlatformMessage, err common_type.PluginError) {
	if err != nil {
		log.ErrorDetails(err)
		return
	}

	ph.OnControlMessage(endpoint, msg)
	ph.OnResourceMessage(endpoint, msg)
}

// OnControlMessage 控制流
func (ph *PlatformHandler) OnControlMessage(endpoint *connect.EndpointInfo, msg *protocol.PlatformMessage) {
	control := msg.GetControl()
	if control == nil {
		return
	}

	if control.GetReport() != nil {
		ph.onHostReport(msg)
	}
}

// OnResourceMessage 资源
func (ph *PlatformHandler) OnResourceMessage(endpoint *connect.EndpointInfo, msg *protocol.PlatformMessage) {
	if msg.GetResource() == nil {
		return
	}

	log.Info("【GET】message.GetResource() GetSeqNo: %d", msg.GetHeader().GetSeqNo())
	resp := resource.NewExecutor(msg).Execute()
	if resp != nil {
		log.Info("【SEND】message.GetResource() GetSeqNo: %d", resp.GetHeader().GetSeqNo())
		if err := ph.SendOnly(resp); err != nil {
			log.ErrorDetails(errors.Trace(err))
		}
	}
}

func (ph *PlatformHandler) onHostReport(msg *protocol.PlatformMessage) {
	report := msg.GetControl().GetReport()
	if report == nil {
		return
	}

	log.Info("【GET】message.HostReport. GetSeqNo: %d", msg.GetHeader().GetSeqNo())

	if logMsg := report.GetLog(); logMsg != nil {
		go func() {
			logger, err := log2.NewLogger(config.StringOrPanic("platform.log_path"))
			if err != nil {
				log.ErrorDetails(err)
				return
			}
			for _, l := range logMsg {
				logger.Log(l)
			}
		}()
	}

	hostID := report.GetHost().GetHostID()
	if hostID == "" {
		return
	}
	_host := host.NewHost(report, common_type.HostStatusNormal)
	ph.pool.Add(_host)
}

// TODO
func (ph *PlatformHandler) StartPlugin(instanceID string) {
	ph.pool.Range(func(hostID string, _host common_type.IHost) bool {
		//if instanceID != hostID {
		//	return true
		//}

		resp := message_utils.BuildP2HDefaultMessage(_host.GetInfo().ID, _host.GetInfo().Name)
		resp.Control.LifeCycleRequest = &protocol.ControlMessage_PluginLifeCycleRequestMessage{
			Instance: &protocol.PluginInstanceDescriptor{
				Application: &protocol.PluginDescriptor{
					ApplicationID:      "lt1ZZuMd",
					Name:               "上传文件的安全提示",
					Language:           "golang",
					LanguageVersion:    message_utils.VersionString2Pb("1.14.0"),
					ApplicationVersion: message_utils.VersionString2Pb("1.0.0"),
					HostVersion:        message_utils.VersionString2Pb("0.2.0"),
					MinSystemVersion:   message_utils.VersionString2Pb("3.2.0"),
				},
				InstanceID: instanceID,
				HostID:     _host.GetInfo().ID,
			},
			Action:     protocol.ControlMessage_Install,
			Reason:     "",
			OldVersion: nil,
		}
		ph.SendAsync(resp, Timeout, func(input, result *protocol.PlatformMessage, err common_type.PluginError) {
			if err == nil {
				return
			}
			if err.Code() != common_type.MsgTimeOut {
				log.Error("Timeout: %+v", resp)
			} else {
				log.ErrorDetails(err)
			}
			ph.pool.Delete(_host) // 过期失活
		})

		return false
	})
}

func (ph *PlatformHandler) Run() common_type.PluginError {
	log.Info("PlatformHandler Run")

	err := ph.GetZmq().Connect()
	if err != nil {
		return err
	}
	go ph.Heartbeat2Host()
	return nil
}
