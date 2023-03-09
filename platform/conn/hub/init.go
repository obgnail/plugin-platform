package hub

import (
	"github.com/obgnail/plugin-platform/platform/conn/hub/ability"
	"github.com/obgnail/plugin-platform/platform/conn/hub/event"
	"github.com/obgnail/plugin-platform/platform/conn/hub/router"
)

// InitHub 插件在整个生命周期需要向platform注册很多东西, 我们将这些存储在内存中:
//   - event: 插件订阅的事件
//   - router: 插件注册的路由
//   - ability: 插件注册的标准能力
func InitHub() error {
	event.InitEvent()
	router.InitRouter()
	ability.InitAbility()
	return nil
}
