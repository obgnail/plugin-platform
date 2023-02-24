package local

import "github.com/obgnail/plugin-platform/common/common_type"

var _ common_type.IResources = (*LocalResource)(nil)

type LocalResource struct {
	plugin common_type.IPlugin

	log       common_type.PluginLogger
	event     common_type.EventPublisher
	space     common_type.Workspace
	sysDBOp   common_type.SysDB
	localDB   common_type.LocalDB
	aPICoreOp common_type.APICore
	network   common_type.Network
	ability   common_type.Ability
}

func New(plugin common_type.IPlugin) *LocalResource {
	l := &LocalResource{
		log:       Logger,
		event:     NewEvent(plugin),
		space:     NewSpace(plugin),
		sysDBOp:   NewSysDB(plugin),
		localDB:   NewLocalDB(plugin),
		aPICoreOp: NewAPICore(plugin),
		network:   NewNetwork(plugin),
		ability:   NewAbility(plugin),
	}
	return l
}

func (r *LocalResource) GetLogger() common_type.PluginLogger {
	return r.log
}

func (r *LocalResource) GetEventPublisher() common_type.EventPublisher {
	return r.event
}

func (r *LocalResource) GetWorkspace() common_type.Workspace {
	return r.space
}

func (r *LocalResource) GetSysDB() common_type.SysDB {
	return r.sysDBOp
}

func (r *LocalResource) GetLocalDB() common_type.LocalDB {
	return r.localDB
}

func (r *LocalResource) GetAPICore() common_type.APICore {
	return r.aPICoreOp
}

func (r *LocalResource) GetOutDoor() common_type.Network {
	return r.network
}

func (r *LocalResource) GetAbility() common_type.Ability {
	return r.ability
}
