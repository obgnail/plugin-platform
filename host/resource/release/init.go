package release

import (
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/host/handler"
	"github.com/obgnail/plugin-platform/host/resource/local"
)

var _ common_type.IResources = (*Resource)(nil)

type Resource struct {
	plugin common_type.IPlugin

	log     common_type.PluginLogger
	event   common_type.EventPublisher
	space   common_type.Workspace
	sysDB   common_type.SysDB
	localDB common_type.LocalDB
	apiCore common_type.APICore
	network common_type.Network
	ability common_type.Ability
}

func NewResource(plugin common_type.IPlugin) *Resource {
	h := handler.Default(false)
	resource := &Resource{
		event:   local.NewEvent(plugin),
		log:     NewLogger(plugin, h),
		space:   NewSpace(plugin, h),
		localDB: NewLocalDB(plugin, h),
		sysDB:   NewSysDB(plugin, h),
		apiCore: NewAPICore(plugin, h),
		network: NewOutdoor(plugin, h),
		ability: local.NewAbility(plugin),
	}
	h.Run()
	return resource
}

func (r *Resource) GetLogger() common_type.PluginLogger {
	return r.log
}

func (r *Resource) GetEventPublisher() common_type.EventPublisher {
	return r.event
}

func (r *Resource) GetWorkspace() common_type.Workspace {
	return r.space
}

func (r *Resource) GetSysDB() common_type.SysDB {
	return r.sysDB
}

func (r *Resource) GetLocalDB() common_type.LocalDB {
	return r.localDB
}

func (r *Resource) GetAPICore() common_type.APICore {
	return r.apiCore
}

func (r *Resource) GetOutDoor() common_type.Network {
	return r.network
}

func (r *Resource) GetAbility() common_type.Ability {
	return r.ability
}
