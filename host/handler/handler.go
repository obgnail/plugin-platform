package handler

import (
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/config"
	"github.com/obgnail/plugin-platform/common/connect"
	"github.com/obgnail/plugin-platform/common/log"
	"github.com/obgnail/plugin-platform/common/protocol"
	"github.com/obgnail/plugin-platform/common/utils/math"
	"github.com/obgnail/plugin-platform/common/utils/message"
	"github.com/obgnail/plugin-platform/host/resource/common"
	"os"
	"time"
)

const (
	defaultTimeoutSec = 30
	RetryReconnectSec = 9
)

var Timeout = time.Duration(config.Int("host.timeout_sec", defaultTimeoutSec)) * time.Second
var RetryReconnectInterval = time.Duration(config.Int("host.retry_reconnect_sec", RetryReconnectSec)) * time.Second

var _ common.Sender = (*HostHandler)(nil)
var _ connect.ConnectionHandler = (*HostHandler)(nil)

type HostHandler struct {
	descriptor   *protocol.HostDescriptor // 存储host的信息
	conn         *connect.Connection      // 负责和platform的通讯
	mounter      *PluginMounter           // 负责挂载插件
	caller       *PluginCaller            // 负责call插件的http
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
			HostVersion:      message.VersionString2Pb(hostVersion),
			MinSystemVersion: message.VersionString2Pb(minSysVersion),
			LanguageVersion:  message.VersionString2Pb(langVersion),
		},
	}

	log.Info("new host: %+v", handler.descriptor)

	zmq := connect.NewZmq(id, name, addr, connect.SocketTypeDealer, connect.RoleHost).SetPacker(&connect.ProtoPacker{})
	handler.conn = connect.NewConnection(zmq, handler)
	handler.mounter = NewMounter(handler, isLocal)
	handler.caller = NewCaller()
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

// InitReport 向platform报告，启动消息循环，等待control指令与其他消息
func (h *HostHandler) InitReport() {
	msg := message.BuildHostReportInitMessage(h.descriptor)
	if err := h.SendOnly(msg); err != nil {
		log.PEDetails(err)
	}
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
	log.PEDetails(err)
	log.Warn("OnError: %s", h.descriptor.Name)
	if err.Code() != common_type.EndpointReceiveErr {
		os.Exit(1)
	}
	// 默认等待3个心跳周期后重新尝试连接,若还连接不上,则退出
	time.Sleep(RetryReconnectInterval)
	if e := h.conn.Connect(); e != nil {
		os.Exit(1)
	}
	h.InitReport()
}

func (h *HostHandler) OnLifeCycle(msg *protocol.PlatformMessage) {
	req := msg.Control.LifeCycleRequest
	host := req.Host
	oldVersion := req.OldVersion
	action := req.Action
	instance := req.Instance
	app := instance.Application
	appID := app.ApplicationID
	appVer := app.ApplicationVersion

	log.Trace("【GET】LifeCycle. [Action]: %d. [appID]: %s. [instanceID]: %s", int32(action), appID, instance)

	resp := &protocol.PlatformMessage{
		Header: &protocol.RouterMessage{
			SeqNo:    msg.Header.SeqNo,
			Source:   msg.Header.Distinct,
			Distinct: msg.Header.Source,
		},
		Control: &protocol.ControlMessage{
			LifeCycleResponse: &protocol.ControlMessage_PluginLifeCycleResponseMessage{
				Host:     host,
				Instance: instance,
				Result:   true, // 这个值后面可能会被修改
				Error:    nil,  // 这个值后面可能会被修改
			},
		},
	}

	// 发送响应数据
	defer func() {
		resp.Control.HostReport = h.buildReportMessage() // 及时报告
		if err := h.SendOnly(resp); err != nil {
			log.PEDetails(err)
			log.Error("appID: %s appVersion: %s", appID, appVer)
		}
	}()

	_plugin, instanceDesc, err := h.mountPlugin(app, instance.InstanceID)
	if err != nil {
		h.whenLifecycleError(resp, action, err)
		return
	}

	err = h.doAction(action, _plugin, oldVersion)
	if err != nil {
		h.whenLifecycleError(resp, action, err)
		return
	}

	h.changePoolStatus(action, instanceDesc)

	return
}

func (h *HostHandler) whenLifecycleError(resp *protocol.PlatformMessage,
	action protocol.ControlMessage_PluginActionType, err common_type.PluginError) {
	log.PEDetails(err)
	resp.Control.LifeCycleResponse.Result = false
	resp.Control.LifeCycleResponse.Error = &protocol.ErrorMessage{
		Code:  int64(h.getLifeCycleErrorCode(action)),
		Error: err.Error(),
		Msg:   err.Msg(),
	}
}

func (h *HostHandler) changePoolStatus(
	action protocol.ControlMessage_PluginActionType,
	instanceDesc common_type.IInstanceDescription,
) {
	switch action {
	case protocol.ControlMessage_UnInstall:
		h.instancePool.DeleteMountedAndRunning(instanceDesc.InstanceID())
	case protocol.ControlMessage_Enable:
		h.instancePool.AddRunning(instanceDesc)
	case protocol.ControlMessage_Disable:
		h.instancePool.DeleteRunning(instanceDesc.InstanceID())
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
	case protocol.ControlMessage_CheckState:
		err = plugin.CheckState()
	case protocol.ControlMessage_CheckCompatibility:
		err = plugin.CheckCompatibility()
	}
	return err
}

func (h *HostHandler) mountPlugin(desc *protocol.PluginDescriptor, instanceID string) (
	common_type.IPlugin, *common_type.MockInstanceDesc, common_type.PluginError) {
	instanceDesc := &common_type.MockInstanceDesc{
		PluginInstanceID: instanceID,
		PluginDescriptor: &common_type.MockPluginDescriptor{
			AppID:      desc.ApplicationID,
			PluginName: desc.Name,
			Lang:       desc.Language,
			LangVer:    message.VersionPb2String(desc.LanguageVersion),
			AppVer:     message.VersionPb2String(desc.ApplicationVersion),
			HostVer:    message.VersionPb2String(desc.HostVersion),
			MinSysVer:  message.VersionPb2String(desc.MinSystemVersion),
		},
	}
	setup, err := h.mounter.Setup(instanceDesc)
	if err != nil {
		log.PEDetails(err)
		return nil, instanceDesc, err
	}

	h.instancePool.AddMounted(instanceDesc.PluginInstanceID, setup)

	return setup, instanceDesc, err
}

func (h *HostHandler) getLifeCycleErrorCode(action protocol.ControlMessage_PluginActionType) int {
	var errCode int
	switch action {
	case protocol.ControlMessage_Install:
		errCode = common_type.OnInstallFailure
	case protocol.ControlMessage_Upgrade:
		errCode = common_type.OnUpgradeFailure
	case protocol.ControlMessage_UnInstall:
		errCode = common_type.OnUnInstallFailure
	case protocol.ControlMessage_Enable:
		errCode = common_type.OnEnableFailure
	case protocol.ControlMessage_Disable:
		errCode = common_type.OnDisEnableFailure
	}
	return errCode
}

func (h *HostHandler) getLifeCycleRequest() common_type.LifeCycleRequest {
	headers := map[string]*common_type.HeaderVal{
		"HostID":           {[]string{h.descriptor.HostID}},
		"HostName":         {[]string{h.descriptor.Name}},
		"HostLanguage":     {[]string{h.descriptor.Language}},
		"LanguageVersion":  {[]string{message.VersionPb2String(h.descriptor.LanguageVersion)}},
		"HostVersion":      {[]string{message.VersionPb2String(h.descriptor.HostVersion)}},
		"MinSystemVersion": {[]string{message.VersionPb2String(h.descriptor.MinSystemVersion)}},
	}

	req := common_type.LifeCycleRequest{Headers: headers}
	return req
}

func (h *HostHandler) buildReportMessage() *protocol.ControlMessage_HostReportMessage {
	var running = make(map[string]*protocol.PluginInstanceDescriptor)
	for _, _running := range h.instancePool.ListRunning() {
		instance := message.BuildInstanceDescriptor(_running, h.descriptor.HostID)
		running[_running.InstanceID()] = instance
	}

	var support = make(map[string]*protocol.PluginInstanceDescriptor)
	for _, _mount := range h.instancePool.ListMounted() {
		description := _mount.GetPluginDescription()
		descriptor := message.BuildInstanceDescriptor(description, h.descriptor.HostID)
		// Q:明明是挂载插件还没有运行,为什么有instanceID? A:由platform生成,接着再传过来
		support[description.InstanceID()] = descriptor
	}

	msg := &protocol.ControlMessage_HostReportMessage{
		Host:          h.descriptor,
		InstanceList:  running,
		SupportedList: support,
	}
	return msg
}

func (h *HostHandler) OnHeartbeat(msg *protocol.PlatformMessage) {
	log.Trace("【GET】Heartbeat. %d", msg.Control.Heartbeat)

	toPlatform := &protocol.PlatformMessage{
		Header: &protocol.RouterMessage{
			SeqNo:    msg.Header.SeqNo,
			Source:   msg.Header.Distinct,
			Distinct: msg.Header.Source,
		},
		Control: &protocol.ControlMessage{
			HostReport: h.buildReportMessage(),
		},
	}

	log.Trace("【SND】Heartbeat. %+v", toPlatform.Control.HostReport)

	if err := h.SendOnly(toPlatform); err != nil {
		log.PEDetails(err)
	}
}

func (h *HostHandler) OnKillSelf(msg *protocol.PlatformMessage) {
	log.Warn("kill handler")
	os.Exit(1)
}

func (h *HostHandler) OnKillPlugin(msg *protocol.PlatformMessage) {
	log.Warn("kill plugin: %+v", msg)

	// go plugin 机制没有卸载功能. 只能在pool中将其删除
	h.instancePool.DeleteMountedAndRunning(msg.Control.KillPlugin.InstanceID)

	resp := &protocol.PlatformMessage{
		Header: &protocol.RouterMessage{
			SeqNo:    msg.Header.SeqNo,
			Source:   msg.Header.Distinct,
			Distinct: msg.Header.Source,
		},
		Control: msg.Control,
	}
	if err := h.SendOnly(resp); err != nil {
		log.PEDetails(err)
	}
}

// 查找正在运行插件
func (h *HostHandler) getRunningInstance(msg *protocol.PlatformMessage) common_type.IInstanceDescription {
	if msg.Plugin.Target == nil {
		return nil
	}
	instanceID := msg.Plugin.Target.InstanceID
	desc, _ := h.instancePool.GetRunning(instanceID)
	return desc
}

func (h *HostHandler) OnPluginHttp(msg *protocol.PlatformMessage) {
	request := msg.Plugin.Http.Request
	if request == nil {
		return
	}

	target := msg.Plugin.Target
	appDesc := target.Application
	instanceID := target.InstanceID
	appID := appDesc.ApplicationID
	appVer := appDesc.ApplicationVersion

	log.Trace("【GET】PluginHttp. [appID]: %s. [instanceID]: %s", appID, instanceID)

	resp := &protocol.PlatformMessage{
		Header: &protocol.RouterMessage{
			SeqNo:    msg.Header.SeqNo,
			Source:   msg.Header.Distinct,
			Distinct: msg.Header.Source,
		},
		Plugin: &protocol.PluginMessage{
			Http: &protocol.HttpContextMessage{
				Response: &protocol.HttpResponseMessage{
					StatusCode: int64(400),
					Headers:    make(map[string]*protocol.HeaderVal),
					Body:       nil,
					Error:      nil,
				},
			},
		},
	}

	defer func() {
		if err := h.SendOnly(resp); err != nil {
			log.PEDetails(err)
			log.Error("appID: %s appVersion: %s", appID, appVer)
		}
	}()

	_plugin, _, err := h.mountPlugin(appDesc, instanceID)
	if err != nil {
		h.whenPluginHttpError(resp, err)
		return
	}

	respMsg, e := h.caller.CallPlugin(_plugin, request)
	if e != nil {
		err = common_type.NewPluginError(common_type.CallPluginHttpFailure, e.Error())
		h.whenPluginHttpError(resp, err)
		return
	}
	resp.Plugin.Http.Response = respMsg
}

func (h *HostHandler) whenPluginHttpError(resp *protocol.PlatformMessage, err common_type.PluginError) {
	log.PEDetails(err)
	resp.Plugin.Http.Response.Error = message.BuildErrorMessage(err)
}

func (h *HostHandler) onErrorTarget(msg *protocol.PlatformMessage) {
	log.Error("错误的插件目标. %+v", msg.Plugin)

	err := common_type.NewPluginError(common_type.GetInstanceFailure, "错误的插件目标")

	resp := &protocol.PlatformMessage{
		Header: &protocol.RouterMessage{
			SeqNo:    msg.Header.SeqNo,
			Source:   msg.Header.Distinct,
			Distinct: msg.Header.Source,
		},
		Plugin: &protocol.PluginMessage{
			Http: &protocol.HttpContextMessage{
				Response: &protocol.HttpResponseMessage{
					Error: message.BuildErrorMessage(err),
				},
			},
		},
	}
	if err = h.SendOnly(resp); err != nil {
		log.PEDetails(err)
	}
}

func (h *HostHandler) OnEvent(msg *protocol.PlatformMessage) {
	event := msg.Plugin.Notification
	if event.Type == "" {
		return
	}

	target := msg.Plugin.Target
	appDesc := target.Application
	instanceID := target.InstanceID
	appID := appDesc.ApplicationID
	appVer := appDesc.ApplicationVersion

	log.Trace("【GET】PluginOnEvent. [Type]:%s [appID]: %s. [instanceID]: %s", event.Type, appID, instanceID)

	resp := &protocol.PlatformMessage{
		Header: &protocol.RouterMessage{
			SeqNo:    msg.Header.SeqNo,
			Source:   msg.Header.Distinct,
			Distinct: msg.Header.Source,
		},
		Plugin: &protocol.PluginMessage{
			Notification: &protocol.NotificationMessage{
				Type:  event.Type,
				Error: nil, // 后续可能会修改此值
			},
		},
	}

	defer func() {
		if err := h.SendOnly(resp); err != nil {
			log.PEDetails(err)
			log.Error("appID: %s appVersion: %s", appID, appVer)
		}
	}()

	_plugin, _, err := h.mountPlugin(appDesc, instanceID)
	if err != nil {
		h.whenPluginOnEventError(resp, err)
		return
	}

	if err = _plugin.OnEvent(event.Type, event.Data); err != nil {
		h.whenPluginOnEventError(resp, err)
		return
	}
}

func (h *HostHandler) whenPluginOnEventError(resp *protocol.PlatformMessage, err common_type.PluginError) {
	log.PEDetails(err)
	resp.Plugin.Notification.Error = message.BuildErrorMessage(err)
}

func (h *HostHandler) OnMsg(endpoint *connect.EndpointInfo, msg *protocol.PlatformMessage, err common_type.PluginError) {
	if err != nil {
		log.PEDetails(err)
		return
	}
	h.OnControlMessage(endpoint, msg)
	h.OnPluginMessage(endpoint, msg)
}

func (h *HostHandler) OnPluginMessage(endpoint *connect.EndpointInfo, msg *protocol.PlatformMessage) {
	pluginMessage := msg.GetPlugin()
	if pluginMessage == nil {
		return
	}

	_plugin := h.getRunningInstance(msg)

	if _plugin == nil {
		h.onErrorTarget(msg)
		return
	}

	// Http请求，使用反射 插件实现的http方法
	if pluginMessage.Http != nil {
		go h.OnPluginHttp(msg)
	}

	// 事件
	if pluginMessage.Notification != nil {
		go h.OnEvent(msg)
	}
}

func (h *HostHandler) OnControlMessage(endpoint *connect.EndpointInfo, msg *protocol.PlatformMessage) {
	control := msg.GetControl()
	if control == nil {
		return
	}

	// 处理HB消息 - 返回应答
	if control.Heartbeat > 0 {
		h.OnHeartbeat(msg)
	}

	// 插件的生命周期管理
	if control.GetLifeCycleRequest() != nil {
		h.OnLifeCycle(msg)
	}

	// kill 自己
	if control.GetKill() != nil {
		h.OnKillSelf(msg)
	}

	// kill 插件
	if control.GetKillPlugin() != nil {
		h.OnKillPlugin(msg)
	}
}

func (h *HostHandler) Send(sender common_type.IPlugin, msg *protocol.PlatformMessage) (*protocol.PlatformMessage, common_type.PluginError) {
	h.fillMsg(sender, msg)
	result, err := h.conn.Send(msg, Timeout)
	return result, err
}

func (h *HostHandler) SendAsync(sender common_type.IPlugin, msg *protocol.PlatformMessage, callback connect.CallBack) {
	h.fillMsg(sender, msg)
	h.conn.SendAsync(msg, Timeout, callback)
}

func (h *HostHandler) SendOnly(msg *protocol.PlatformMessage) (err common_type.PluginError) {
	h.fillMsg(nil, msg)
	return h.conn.SendOnly(msg)
}

// fillMsg 添加路由信息
func (h *HostHandler) fillMsg(sender common_type.IPlugin, msg *protocol.PlatformMessage) {
	if msg == nil {
		msg = message.GetInitMessage(nil, nil)
	}
	msg.Header.Source = message.GetHostInfo(h.descriptor.HostID, h.descriptor.Name)
	msg.Header.Distinct = message.GetPlatformInfo()
	if msg.Header.SeqNo == 0 {
		msg.Header.SeqNo = math.CreateCaptcha()
	}
	if msg.Resource != nil && sender != nil {
		msg.Resource.Sender = message.BuildInstanceDescriptor(sender.GetPluginDescription(), h.descriptor.HostID)
		msg.Resource.Host = h.descriptor
	}
}

func (h *HostHandler) Run() common_type.PluginError {
	if err := h.conn.Connect(); err != nil {
		return err
	}
	go func() {
		time.Sleep(time.Second * 1)
		h.InitReport()
	}()
	return nil
}
