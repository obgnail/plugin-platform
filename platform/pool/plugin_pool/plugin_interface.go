package plugin_pool

const (
	PluginStatusUploaded    = 1 // 已上传
	PluginStatusRunning     = 2 // 已启用
	PluginStatusStopping    = 3 // 已安装/已停用
	PluginStatusUnavailable = 4 // 不可用
)

type PluginInterface interface {
	Configuration() *PluginConfig
	Enable() error
	Disable() error
	PluginStatus() int
	SetSubscribe(subscribeRules map[string][]string)
	GetSubscribe() map[string][]string
	PluginInstanceId() string
}

type PluginProcess struct {
	*PluginConfig
	Subscribe map[string][]string
}

func NewPluginProcess(pluginConfig *PluginConfig) *PluginProcess {
	var r = new(PluginProcess)
	r.PluginConfig = pluginConfig
	return r
}

func (p *PluginProcess) Configuration() *PluginConfig {
	return p.PluginConfig
}

func (p *PluginProcess) Enable() error {
	p.Status = PluginStatusRunning
	return nil
}

func (p *PluginProcess) Disable() error {
	p.Status = PluginStatusStopping
	return nil
}

func (p *PluginProcess) PluginStatus() int {
	return p.Status
}

func (p *PluginProcess) PluginInstanceId() string {
	return p.InstanceUUID
}

func (p *PluginProcess) SetSubscribe(subscribeRules map[string][]string) {
	p.Subscribe = subscribeRules
}

func (p *PluginProcess) GetSubscribe() map[string][]string {
	return p.Subscribe
}
