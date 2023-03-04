package common_type

type HostStatus int

const (
	HostStatusNormal HostStatus = 1
	HostStatusDrift  HostStatus = 2
)

type IHost interface {
	//Store()
	GetInfo() HostInfo
	GetStatus() HostStatus
	//KillSelf()
	//KillPlugin(instanceID string)
	//CallPlugin(ctx *RequestContext) (resp *LifecycleResponse, err error)
}

type HostInfo struct {
	ID               string
	Name             string
	Version          string
	MinSystemVersion string
	Language         string
	LanguageVersion  string
	RunningPlugins   map[string]IInstanceDescription
	SupportPlugins   map[string]IInstanceDescription
}
