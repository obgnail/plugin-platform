package handler

import (
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/config"
	"github.com/obgnail/plugin-platform/common/connect"
	"github.com/obgnail/plugin-platform/common/log"
	"github.com/obgnail/plugin-platform/common/message_utils"
	"github.com/obgnail/plugin-platform/common/protocol"
	"os"
	"time"
)

// TODO
var defaultTimeout = 30 * time.Second

var _ connect.FurtherHandler = (*HostHandler)(nil)

type HostHandler struct {
	descriptor   *protocol.HostDescriptor // 存储host的信息
	conn         *connect.BaseHandler     // 负责和platform的通讯
	mounter      *PluginMounter           // 负责挂载插件
	instancePool *InstancePool            // 存储已经挂载的插件
	isLocal      bool                     // host运行在测试环境/生产环境
}

func New(id, name, addr, lang, hostVersion, minSysVersion, langVersion string, isLocal bool) *HostHandler {
	handler := &HostHandler{
		instancePool: &InstancePool{},
		descriptor: &protocol.HostDescriptor{
			HostID:           id,
			Name:             name,
			Language:         lang,
			HostVersion:      message_utils.SplitVersion(hostVersion),
			MinSystemVersion: message_utils.SplitVersion(minSysVersion),
			LanguageVersion:  message_utils.SplitVersion(langVersion),
		},
	}

	zmq := connect.NewZmq(id, name, addr, connect.SocketTypeDealer, connect.RoleHost).SetPacker(&connect.ProtoPacker{})
	handler.conn = connect.NewBaseHandler(zmq, handler)
	handler.mounter = NewMounter(handler, isLocal)
	return handler
}

func Default(isLocal bool) *HostHandler {
	id := config.StringOrPanic("host.id")
	name := config.StringOrPanic("host.name")
	addr := config.StringOrPanic("host.platform_address")
	lang := config.StringOrPanic("host.language")
	hostVersion := config.StringOrPanic("host.version")
	minSysVersion := config.StringOrPanic("host.min_system_version")
	langVersion := config.StringOrPanic("host.language_version")

	h := New(id, name, addr, lang, hostVersion, minSysVersion, langVersion, isLocal)
	return h
}

func (h *HostHandler) GetDescriptor() *protocol.HostDescriptor {
	return h.descriptor
}

func (h *HostHandler) OnConnect() common_type.PluginError {
	log.Info("OnConnect: %s", h.descriptor.Name)
	return nil
}

func (h *HostHandler) OnDisconnect() common_type.PluginError {
	log.Info("OnDisconnect: %s", h.descriptor.Name)
	return nil
}

func (h *HostHandler) OnError(err common_type.PluginError) {
	log.Warn("OnError: %s", h.descriptor.Name)
	if err.Code() != common_type.EndpointReceiveErr {
		os.Exit(1)
	}
	time.Sleep(time.Second * 9)
	if e := h.conn.Connect(); e != nil {
		os.Exit(1)
	}
}

func (h *HostHandler) OnKill(message *protocol.PlatformMessage) {
	log.Info("kill handler")
	os.Exit(1)
}

func (h *HostHandler) OnLifeCycle(msg *protocol.PlatformMessage) {
	resp := &protocol.PlatformMessage{
		Header: &protocol.RouterMessage{
			SeqNo:    msg.Header.SeqNo,
			Source:   msg.Header.Distinct,
			Distinct: msg.Header.Source,
		},
		Control: &protocol.ControlMessage{
			LifeCycleResponse: &protocol.ControlMessage_PluginLifeCycleResponseMessage{
				Host:     msg.Control.LifeCycleRequest.Host,
				Instance: msg.Control.LifeCycleRequest.Instance,
				Result:   true, // 这个值后面可能会被setLifeCycleRespError()修改
				Error:    nil,  // 这个值后面可能会被setLifeCycleRespError()修改
			},
		},
	}

	// 发送响应数据
	defer func() {
		if err := h.conn.SendOnly(resp); err != nil {
			log.Error("appId: %s appVersion: %s hh.SendMessage err: %s",
				msg.Control.LifeCycleRequest.Instance.Application.ApplicationID,
				msg.Control.LifeCycleRequest.Instance.Application.ApplicationVersion, err.Error())
		}
	}()

	instance := msg.Control.LifeCycleRequest.Instance
	action := msg.Control.LifeCycleRequest.Action
	oldVersion := msg.Control.LifeCycleRequest.OldVersion
	app := instance.Application

	instanceDesc := &MockInstanceDesc{
		instanceID: instance.InstanceID,
		pluginDescriptor: &MockPluginDescriptor{
			appID:     app.ApplicationID,
			name:      app.Name,
			lang:      app.Language,
			langVer:   message_utils.GetVersionString(app.LanguageVersion),
			appVer:    message_utils.GetVersionString(app.ApplicationVersion),
			hostVer:   message_utils.GetVersionString(app.HostVersion),
			minSysVer: message_utils.GetVersionString(app.MinSystemVersion),
		},
	}

	_plugin, err := h.mountPlugin(instanceDesc)
	if err != nil {
		h.setLifeCycleRespError(resp, action, err)
		return
	}

	err = h.doAction(action, _plugin, oldVersion)
	if err != nil {
		h.setLifeCycleRespError(resp, action, err)
		return
	}

	h.changePoolStatus(action, instanceDesc)

	return
}

func (h *HostHandler) changePoolStatus(
	action protocol.ControlMessage_PluginActionType,
	instanceDesc common_type.IInstanceDescription,
) {
	switch action {
	case protocol.ControlMessage_UnInstall:
		h.instancePool.DeleteInstance(instanceDesc.InstanceID())
	case protocol.ControlMessage_Start:
		h.instancePool.StartInstance(instanceDesc)
	case protocol.ControlMessage_Disable:
		h.instancePool.StopInstance(instanceDesc.InstanceID())
	}
}

func (h *HostHandler) doAction(
	action protocol.ControlMessage_PluginActionType,
	plugin common_type.IPlugin,
	oldVersion *protocol.PluginDescriptor,
) (err common_type.PluginError) {
	request := h.getLifeCycleRequest()

	switch action {
	case protocol.ControlMessage_Install:
		err = plugin.Install(request)
	case protocol.ControlMessage_UnInstall:
		err = plugin.UnInstall(request)
	case protocol.ControlMessage_Enable:
		err = plugin.Enable(request)
	case protocol.ControlMessage_Disable:
		err = plugin.Disable(request)
	case protocol.ControlMessage_Upgrade:
		major := oldVersion.ApplicationVersion.GetMajor()
		minor := oldVersion.ApplicationVersion.GetMinor()
		revision := oldVersion.ApplicationVersion.GetRevision()
		ver := common_type.NewVersion(int(major), int(minor), int(revision))
		err = plugin.Upgrade(ver, request)
	case protocol.ControlMessage_Start:
		err = plugin.Start(request)
	case protocol.ControlMessage_Stop:
		err = plugin.Stop(request)
	case protocol.ControlMessage_CheckState:
		err = plugin.CheckState()
	case protocol.ControlMessage_CheckCompatibility:
		err = plugin.CheckCompatibility()
	}
	return err
}

func (h *HostHandler) mountPlugin(instanceDesc *MockInstanceDesc) (common_type.IPlugin, common_type.PluginError) {
	_plugin, _, _ := h.instancePool.GetPlugin(instanceDesc.InstanceID())

	setupPlugin, err := h.mounter.Mount(_plugin, instanceDesc)
	if err != nil {
		return nil, err
	}

	h.instancePool.Add(instanceDesc.instanceID, setupPlugin) // 此插件已经成功挂载
	return setupPlugin, nil
}

func (h *HostHandler) setLifeCycleRespError(resp *protocol.PlatformMessage,
	action protocol.ControlMessage_PluginActionType, err common_type.PluginError) {
	log.ErrorDetails(err)
	resp.Control.LifeCycleResponse.Result = false
	resp.Control.LifeCycleResponse.Error = &protocol.ErrorMessage{
		Code:  int64(h.getLifeCycleErrorCode(action)),
		Error: err.Error(),
		Msg:   err.Msg(),
	}
}

func (h *HostHandler) getLifeCycleErrorCode(action protocol.ControlMessage_PluginActionType) int {
	var errcode int
	switch action {
	case protocol.ControlMessage_Install:
		errcode = common_type.OnInstallFailure
	case protocol.ControlMessage_Upgrade:
		errcode = common_type.OnUpgradeFailure
	case protocol.ControlMessage_UnInstall:
		errcode = common_type.OnUnInstallFailure
	case protocol.ControlMessage_Enable:
		errcode = common_type.OnEnableFailure
	case protocol.ControlMessage_Disable:
		errcode = common_type.OnDisEnableFailure
	case protocol.ControlMessage_Start:
		errcode = common_type.OnstartFailure
	case protocol.ControlMessage_Stop:
		errcode = common_type.OnstopFailure
	}
	return errcode
}

func (h *HostHandler) getLifeCycleRequest() common_type.LifeCycleRequest {
	headers := map[string]*common_type.HeaderVal{
		"HostID":           {[]string{h.descriptor.HostID}},
		"HostName":         {[]string{h.descriptor.Name}},
		"HostLanguage":     {[]string{h.descriptor.Language}},
		"LanguageVersion":  {[]string{message_utils.GetVersionString(h.descriptor.LanguageVersion)}},
		"HostVersion":      {[]string{message_utils.GetVersionString(h.descriptor.HostVersion)}},
		"MinSystemVersion": {[]string{message_utils.GetVersionString(h.descriptor.MinSystemVersion)}},
	}

	req := common_type.LifeCycleRequest{Headers: headers}
	return req
}

func (h *HostHandler) OnMsg(endpoint *connect.EndpointInfo, msg *protocol.PlatformMessage, err common_type.PluginError) {
	if err != nil {
		log.ErrorDetails(err)
		return
	}

	control := msg.GetControl()
	if control != nil {
		// 插件的生命周期管理
		if control.GetLifeCycleRequest() != nil {
			h.OnLifeCycle(msg)
		}
		// kill 自己
		if control.GetKill() != nil {
			h.OnKill(msg)
		}
	}

	// 资源请求的应答
	if msg.GetResource() != nil {
		log.Info("%+v", msg)
	}
}

func (h *HostHandler) Send(sender common_type.IPlugin, msg *protocol.PlatformMessage) (*protocol.PlatformMessage, common_type.PluginError) {
	h.fillMsg(sender, msg)
	result, err := h.conn.Send(msg, defaultTimeout)
	return result, err
}

func (h *HostHandler) SendAsync(sender common_type.IPlugin, msg *protocol.PlatformMessage, callback connect.CallBack) {
	h.fillMsg(sender, msg)
	h.conn.SendAsync(msg, defaultTimeout, callback)
}

// fillMsg 添加路由信息
func (h *HostHandler) fillMsg(sender common_type.IPlugin, msg *protocol.PlatformMessage) {
	if msg.GetResource() != nil {
		msg.Resource.Sender = message_utils.BuildInstanceDescriptor(sender, h.descriptor.HostID)
		msg.Resource.Host = h.descriptor
		////补全日志信息
		//if msg.GetResource().GetLog() != nil {
		//	msg.Resource.Log.PluginInstanceDescriptor = pluginInstanceDescriptor
		//	msg.Resource.Log.HostDescriptor = hh.hostDescriptor
		//}
	}
}

func (h *HostHandler) Run() common_type.PluginError {
	return h.conn.Connect()
}
