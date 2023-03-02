package handler

import (
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/config"
	"github.com/obgnail/plugin-platform/common/connect"
	"github.com/obgnail/plugin-platform/common/log"
	"github.com/obgnail/plugin-platform/common/math"
	"github.com/obgnail/plugin-platform/common/message_utils"
	"github.com/obgnail/plugin-platform/common/protocol"
	host_handler "github.com/obgnail/plugin-platform/host/handler"
	"os"
	"time"
)

var Timeout = time.Duration(config.Int("host_boot.timeout_sec", 30)) * time.Second

type HostBootHandler struct {
	descriptor *protocol.HostBootDescriptor
	conn       *connect.BaseHandler // 负责和platform的通讯
}

func New(id, name, addr, version string) *HostBootHandler {
	h := &HostBootHandler{
		descriptor: &protocol.HostBootDescriptor{
			BootID:      id,
			Name:        name,
			BootVersion: message_utils.VersionString2Pb(version),
		},
	}
	zmq := connect.NewZmq(id, name, addr, connect.SocketTypeDealer, connect.RoleHost).SetPacker(&connect.ProtoPacker{})
	h.conn = connect.NewBaseHandler(zmq, h)
	return h
}

func Default() *HostBootHandler {
	id := config.StringOrPanic("host_boot.id")
	name := config.StringOrPanic("host_boot.name")
	version := config.StringOrPanic("host_boot.version")
	addr := config.StringOrPanic("host_boot.platform_address")
	return New(id, name, addr, version)
}

func (h *HostBootHandler) OnConnect() common_type.PluginError {
	log.Info("host_boot OnConnect")
	return nil
}

func (h *HostBootHandler) OnDisconnect() common_type.PluginError {
	log.Info("host_boot OnDisconnect")
	return nil
}

func (h *HostBootHandler) OnError(err common_type.PluginError) {
	log.Warn("OnError %+v", err)
	if err.Code() != common_type.EndpointReceiveErr {
		os.Exit(1)
	}
	time.Sleep(time.Second * 9)
	if e := h.conn.Connect(); e != nil {
		os.Exit(1)
	}
}

func (h *HostBootHandler) OnMsg(endpoint *connect.EndpointInfo, msg *protocol.PlatformMessage, err common_type.PluginError) {
	if err != nil {
		log.ErrorDetails(err)
		return
	}

	control := msg.GetControl()
	if control == nil {
		return
	}

	// 处理HB消息 - 返回应答
	if control.Heartbeat > 0 {
		h.OnHeartbeat(msg)
	}

	if control.StartHost != nil {
		h.OnStartHost(msg)
	}
}

func (h *HostBootHandler) OnStartHost(msg *protocol.PlatformMessage) {
	host := h.newHost(msg)

	resp := &protocol.PlatformMessage{
		Header: &protocol.RouterMessage{
			SeqNo:    msg.Header.SeqNo,
			Source:   msg.Header.Distinct,
			Distinct: msg.Header.Source,
		},
		Control: &protocol.ControlMessage{
			StartHost: &protocol.ControlMessage_StartHostMessage{
				Host:   host.GetDescriptor(),
				Result: true,
				Error:  nil,
			},
		},
	}
	log.Info("start host: %+v", resp)
	if err := h.SendOnly(resp); err != nil {
		log.ErrorDetails(err)
	}
}

func (h *HostBootHandler) newHost(msg *protocol.PlatformMessage) *host_handler.HostHandler {
	host := msg.Control.StartHost.Host
	id := host.HostID
	name := host.Name

	addr := config.StringOrPanic("host.platform_address")
	lang := config.StringOrPanic("host.language")
	hostVersion := config.StringOrPanic("host.version")
	minSysVersion := config.StringOrPanic("host.min_system_version")
	langVersion := config.StringOrPanic("host.language_version")

	hostHandler := host_handler.New(id, name, addr, lang, hostVersion, minSysVersion, langVersion, false)
	go hostHandler.Run()
	return hostHandler
}

func (h *HostBootHandler) OnHeartbeat(msg *protocol.PlatformMessage) {
	h.Report()
}

// Report 向platform报告，启动消息循环，等待control指令与其他消息
func (h *HostBootHandler) Report() {
	msg := message_utils.BuildHostBootReportInitMessage(h.descriptor)
	if err := h.SendOnly(msg); err != nil {
		log.ErrorDetails(err)
	}
}

func (h *HostBootHandler) Send(msg *protocol.PlatformMessage) (*protocol.PlatformMessage, common_type.PluginError) {
	h.fillMsg(msg)
	result, err := h.conn.Send(msg, Timeout)
	return result, err
}

func (h *HostBootHandler) SendAsync(msg *protocol.PlatformMessage, callback connect.CallBack) {
	h.fillMsg(msg)
	h.conn.SendAsync(msg, Timeout, callback)
}

func (h *HostBootHandler) SendOnly(msg *protocol.PlatformMessage) (err common_type.PluginError) {
	h.fillMsg(msg)
	return h.conn.SendOnly(msg)
}

// fillMsg 添加路由信息
func (h *HostBootHandler) fillMsg(msg *protocol.PlatformMessage) {
	if msg == nil {
		msg = message_utils.GetInitMessage(nil, nil)
	}
	msg.Header.Source = message_utils.GetHostBootInfo(h.descriptor.BootID, h.descriptor.Name)
	msg.Header.Distinct = message_utils.GetPlatformInfo()
	if msg.Header.SeqNo == 0 {
		msg.Header.SeqNo = math.CreateCaptcha()
	}
}

func (h *HostBootHandler) Run() common_type.PluginError {
	if err := h.conn.Connect(); err != nil {
		return err
	}
	go func() {
		time.Sleep(time.Second * 1)
		h.Report()
	}()
	return nil
}
