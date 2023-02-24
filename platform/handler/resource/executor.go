package resource

import (
	"github.com/obgnail/plugin-platform/common/message_utils"
	"github.com/obgnail/plugin-platform/common/protocol"
	"github.com/obgnail/plugin-platform/platform/handler/resource/db"
	"github.com/obgnail/plugin-platform/platform/handler/resource/network"
	"github.com/obgnail/plugin-platform/platform/handler/resource/work_space"
)

type Executor struct {
	Request  *protocol.PlatformMessage
	Response *protocol.PlatformMessage // 如果需要返回信息,此值不为空
}

func NewExecutor(request *protocol.PlatformMessage) *Executor {
	resource := &Executor{
		Request:  request,
		Response: message_utils.GetResourceInitMessage(request),
	}
	return resource
}

// Execute 注意:一条请求可能包含多种资源操作
func (r *Executor) Execute() (resp *protocol.PlatformMessage) {
	// workspace
	if r.Request.GetResource().GetWorkspace() != nil {
		work_space.NewWorkSpace(r.Request, r.Response).Execute()
	}
	// localDB、sysDB
	if r.Request.GetResource().GetDatabase() != nil {
		db.NewDataBase(r.Request, r.Response).Execute()
	}
	// apiCore、outdoor
	if r.Request.GetResource().GetHttp() != nil {
		network.NewNetWork(r.Request, r.Response).Execute()
	}

	return r.Response
}
