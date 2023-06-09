package common_type

type IPlugin interface {
	// 程序实现
	Assign(pid IInstanceDescription, resources IResources) PluginError
	GetPluginDescription() IInstanceDescription

	// 业务动作
	Enable(LifeCycleRequest) PluginError
	Disable(LifeCycleRequest) PluginError
	Install(LifeCycleRequest) PluginError
	UnInstall(LifeCycleRequest) PluginError
	Upgrade(IVersion, LifeCycleRequest) PluginError
	CheckState() PluginError
	CheckCompatibility() PluginError

	// 事件相关
	OnEvent(eventType string, payload []byte) PluginError

	// 配置修改
	OnConfigChange(configKey string, originValue, newValue []string) PluginError

	// 外部请求
	OnExternalHttpRequest(request *HttpRequest) *HttpResponse
}

type IInstanceDescription interface {
	PluginDescription() IPluginDescriptor
	InstanceID() string
}

type IPluginDescriptor interface {
	ApplicationID() string
	Name() string
	Language() string
	LanguageVersion() IVersion
	ApplicationVersion() IVersion
	HostVersion() IVersion
	MinSystemVersion() IVersion
}

type IResources interface {
	GetLogger() PluginLogger
	GetWorkspace() Workspace
	GetLocalDB() LocalDB
	GetEventPublisher() EventPublisher
	GetSysDB() SysDB
	GetAPICore() APICore
	GetOutDoor() Network
	GetAbility() Ability
}

type LifeCycleRequest struct {
	Headers map[string]*HeaderVal
}

type HeaderVal struct {
	Val []string
}

type PluginFactory interface {
	GetPlugin() IPlugin
}
