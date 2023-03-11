package router

import (
	"github.com/obgnail/plugin-platform/platform/conn/hub/router/http_router"
	"github.com/obgnail/plugin-platform/platform/service/common"
	"sync"
)

var mu sync.Mutex
var r *PluginRouter

func InitRouter() {
	mu.Lock()
	defer mu.Unlock()
	r = NewRouter()
}

func RegisterRouter(instanceID string, apis []*common.Api) error {
	mu.Lock()
	defer mu.Unlock()
	return r.Register(instanceID, apis)
}

func MatchRouter(Type, method, url string) *http_router.RouterInfo {
	mu.Lock()
	defer mu.Unlock()
	return r.Match(Type, method, url)
}
