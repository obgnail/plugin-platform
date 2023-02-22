package local

import common "github.com/obgnail/plugin-platform/common_type"

var _ common.IResources = (*LocalResource)(nil)

type LocalResource struct {
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

func New(plugin common.IPlugin) *LocalResource {
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

func (r *LocalResource) GetLogger() common.PluginLogger {
	return r.log
}

func (r *LocalResource) GetEventPublisher() common.EventPublisher {
	return r.event
}

func (r *LocalResource) GetWorkspace() common.Workspace {
	return r.space
}

func (r *LocalResource) GetSysDB() common.SysDB {
	return r.sysDBOp
}

func (r *LocalResource) GetLocalDB() common.LocalDB {
	return r.localDB
}

func (r *LocalResource) GetAPICore() common.APICore {
	return r.aPICoreOp
}

func (r *LocalResource) GetOutDoor() common.Network {
	return r.network
}

func (r *LocalResource) GetAbility() common.Ability {
	return r.ability
}
