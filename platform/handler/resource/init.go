package resource

import (
	"github.com/obgnail/plugin-platform/common/message_utils"
	"github.com/obgnail/plugin-platform/common/protocol"
	"github.com/obgnail/plugin-platform/platform/handler/resource/db"
	"github.com/obgnail/plugin-platform/platform/handler/resource/work_space"
)

type Executor struct {
	Source   *protocol.PlatformMessage
	Distinct *protocol.PlatformMessage
}

func NewExecutor(source *protocol.PlatformMessage) *Executor {
	resource := &Executor{
		Source:   source,
		Distinct: message_utils.GetResourceInitMessage(source),
	}
	return resource
}

func (r *Executor) Execute() {
	// 文件操作
	if r.Source.GetResource().GetWorkspace() != nil {
		work_space.NewWorkSpace(r.Source, r.Distinct).Execute()
	}
	//
	if r.Source.GetResource().GetDatabase() != nil {
		db.NewDataBase(r.Source, r.Distinct).Execute()
	}
}
