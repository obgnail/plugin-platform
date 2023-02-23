package resource

import (
	"github.com/BangWork/ones-platform-api/protocol/build_message"
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/protocol"
)

var _ common_type.IResources = (*Resource)(nil)

type Resource struct {
	SourceMessage   *protocol.PlatformMessage
	DistinctMessage *protocol.PlatformMessage

	log       common_type.PluginLogger
	event     common_type.EventPublisher
	space     common_type.Workspace
	sysDBOp   common_type.SysDB
	localDB   common_type.LocalDB
	aPICoreOp common_type.APICore
	network   common_type.Network
	ability   common_type.Ability
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
	return r.sysDBOp
}

func (r *Resource) GetLocalDB() common_type.LocalDB {
	return r.localDB
}

func (r *Resource) GetAPICore() common_type.APICore {
	return r.aPICoreOp
}

func (r *Resource) GetOutDoor() common_type.Network {
	return r.network
}

func (r *Resource) GetAbility() common_type.Ability {
	return r.ability
}
