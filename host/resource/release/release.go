package release

import (
	common "github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/host/handler"
	"github.com/obgnail/plugin-platform/host/resource/local"
)

var _ common.IResources = (*ReleaseResource)(nil)

type ReleaseResource struct {
	*handler.HostHandler

	plugin common.IPlugin

	log       common.PluginLogger
	event     common.EventPublisher
	space     common.Workspace
	sysDBOp   common.SysDB
	localDB   common.LocalDB
	aPICoreOp common.APICore
	network   common.Network
	ability   common.Ability
}

func NewReleaseResource(plugin common.IPlugin) *ReleaseResource {
	h := handler.Default(false)
	resource := &ReleaseResource{
		HostHandler: h,
		log:         local.Logger,
		event:       local.NewEvent(plugin),
		space:       NewSpace(plugin, h),
		sysDBOp:     local.NewSysDB(plugin),
		localDB:     local.NewLocalDB(plugin),
		aPICoreOp:   local.NewAPICore(plugin),
		network:     local.NewNetwork(plugin),
		ability:     local.NewAbility(plugin),
	}
	resource.Run()
	return resource
}

func (r *ReleaseResource) GetLogger() common.PluginLogger {
	return r.log
}

func (r *ReleaseResource) GetEventPublisher() common.EventPublisher {
	return r.event
}

func (r *ReleaseResource) GetWorkspace() common.Workspace {
	return r.space
}

func (r *ReleaseResource) GetSysDB() common.SysDB {
	return r.sysDBOp
}

func (r *ReleaseResource) GetLocalDB() common.LocalDB {
	return r.localDB
}

func (r *ReleaseResource) GetAPICore() common.APICore {
	return r.aPICoreOp
}

func (r *ReleaseResource) GetOutDoor() common.Network {
	return r.network
}

func (r *ReleaseResource) GetAbility() common.Ability {
	return r.ability
}
