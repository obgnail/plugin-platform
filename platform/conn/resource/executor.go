package resource

import (
	"github.com/obgnail/plugin-platform/common/protocol"
	"github.com/obgnail/plugin-platform/common/utils/message"
	"github.com/obgnail/plugin-platform/platform/conn/resource/ability"
	"github.com/obgnail/plugin-platform/platform/conn/resource/db"
	"github.com/obgnail/plugin-platform/platform/conn/resource/event_publisher"
	"github.com/obgnail/plugin-platform/platform/conn/resource/log"
	"github.com/obgnail/plugin-platform/platform/conn/resource/network"
	"github.com/obgnail/plugin-platform/platform/conn/resource/work_space"
)

type Executor struct {
	Request  *protocol.PlatformMessage
	Response *protocol.PlatformMessage // 如果需要返回信息,此值不为空
}

func NewExecutor(request *protocol.PlatformMessage) *Executor {
	resource := &Executor{
		Request:  request,
		Response: message.GetResourceInitMessage(request),
	}
	return resource
}

// Execute 注意:一条请求可能包含多种资源操作
func (r *Executor) Execute() (resp *protocol.PlatformMessage) {
	resource := r.Request.GetResource()

	// log
	if resource.GetLog() != nil {
		log.NewLog(r.Request, r.Response).Execute()
	}
	// workspace
	if resource.GetWorkspace() != nil {
		work_space.NewWorkSpace(r.Request, r.Response).Execute()
	}
	// localDB、sysDB
	if resource.GetDatabase() != nil {
		db.NewDataBase(r.Request, r.Response).Execute()
	}
	// apiCore、outdoor
	if resource.GetHttp() != nil {
		network.NewNetWork(r.Request, r.Response).Execute()
	}
	// event
	if resource.GetEvent() != nil {
		event_publisher.NewEvent(r.Request, r.Response).Execute()
	}
	// ability
	if resource.GetAbility() != nil {
		ability.NewAbility(r.Request, r.Response).Execute()
	}

	return r.Response
}
