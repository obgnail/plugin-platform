package router

import (
	"github.com/obgnail/plugin-platform/platform/conn/hub/router/http_router"
	"github.com/obgnail/plugin-platform/platform/service/types"
	"sync"
)

var mu sync.Mutex
var r *PluginRouter

func InitRouter() {
	mu.Lock()
	defer mu.Unlock()
	r = NewRouter()
}

func RegisterRouter(instanceID string, apis []*types.Api) error {
	mu.Lock()
	defer mu.Unlock()
	return r.Register(instanceID, apis)
}

func DeleteRouter(instanceID string) {
	mu.Lock()
	defer mu.Unlock()
	r.Delete(instanceID)
}

func MatchRouter(Type, method, url string) *http_router.RouterInfo {
	mu.Lock()
	defer mu.Unlock()
	return r.Match(Type, method, url)
}
