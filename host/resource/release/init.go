package release

import (
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/host/resource/common"
)

var _ common_type.IResources = (*Resource)(nil)

type Resource struct {
	Plugin common_type.IPlugin
	Sender common.Sender

	log     common_type.PluginLogger
	event   common_type.EventPublisher
	space   common_type.Workspace
	sysDB   common_type.SysDB
	localDB common_type.LocalDB
	apiCore common_type.APICore
	network common_type.Network
	ability common_type.Ability
}

func (r *Resource) GetLogger() common_type.PluginLogger {
	if r.log == nil {
		r.log = NewLogger(r.Plugin, r.Sender)
	}
	return r.log
}

func (r *Resource) GetEventPublisher() common_type.EventPublisher {
	if r.event == nil {
		r.event = NewEvent(r.Plugin, r.Sender)
	}
	return r.event
}

func (r *Resource) GetWorkspace() common_type.Workspace {
	if r.space == nil {
		r.space = NewSpace(r.Plugin, r.Sender)
	}
	return r.space
}

func (r *Resource) GetSysDB() common_type.SysDB {
	if r.sysDB == nil {
		r.sysDB = NewSysDB(r.Plugin, r.Sender)
	}
	return r.sysDB
}

func (r *Resource) GetLocalDB() common_type.LocalDB {
	if r.localDB == nil {
		r.localDB = NewLocalDB(r.Plugin, r.Sender)
	}
	return r.localDB
}

func (r *Resource) GetAPICore() common_type.APICore {
	if r.apiCore == nil {
		r.apiCore = NewAPICore(r.Plugin, r.Sender)
	}
	return r.apiCore
}

func (r *Resource) GetOutDoor() common_type.Network {
	if r.network == nil {
		r.network = NewOutdoor(r.Plugin, r.Sender)
	}
	return r.network
}

func (r *Resource) GetAbility() common_type.Ability {
	if r.ability == nil {
		r.ability = NewAbility(r.Plugin, r.Sender)
	}
	return r.ability
}
