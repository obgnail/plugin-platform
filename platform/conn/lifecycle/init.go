package lifecycle

import (
	"github.com/obgnail/plugin-platform/platform/conn/lifecycle/event"
	"github.com/obgnail/plugin-platform/platform/conn/lifecycle/router"
)

func InitLifecycleTool() error {
	event.InitEvent()
	router.InitRouter()
	return nil
}
