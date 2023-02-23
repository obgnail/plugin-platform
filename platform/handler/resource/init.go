package resource

import (
	"github.com/BangWork/ones-platform-api/protocol/build_message"
	common "github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/protocol"
)

var _ common.IResources = (*Resource)(nil)

type Resource struct {
	SourceMessage   *protocol.PlatformMessage
	DistinctMessage *protocol.PlatformMessage

	log       common.PluginLogger
	event     common.EventPublisher
	space     common.Workspace
	sysDBOp   common.SysDB
	localDB   common.LocalDB
	aPICoreOp common.APICore
	network   common.Network
	ability   common.Ability
}

func New(sourceMessage *protocol.PlatformMessage) *Resource {
	resource := &Resource{
		SourceMessage:   sourceMessage,
		DistinctMessage: build_message.GetResourceInitMessage(),
		//log:       Logger,
		//event:     NewEvent(plugin),
		space: NewSpace(plugin),
		//sysDBOp:   NewSysDB(plugin),
		//localDB:   NewLocalDB(plugin),
		//aPICoreOp: NewAPICore(plugin),
		//network:   NewNetwork(plugin),
		//ability:   NewAbility(plugin),
	}
	return resource
}

func (r *Resource) GetLogger() common.PluginLogger {
	return r.log
}

func (r *Resource) GetEventPublisher() common.EventPublisher {
	return r.event
}

func (r *Resource) GetWorkspace() common.Workspace {
	return r.space
}

func (r *Resource) GetSysDB() common.SysDB {
	return r.sysDBOp
}

func (r *Resource) GetLocalDB() common.LocalDB {
	return r.localDB
}

func (r *Resource) GetAPICore() common.APICore {
	return r.aPICoreOp
}

func (r *Resource) GetOutDoor() common.Network {
	return r.network
}

func (r *Resource) GetAbility() common.Ability {
	return r.ability
}
