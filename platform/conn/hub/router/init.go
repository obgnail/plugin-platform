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

func RegisterRouter(apis []*common.Api, instanceID string) error {
	mu.Lock()
	defer mu.Unlock()
	return r.Add(apis, instanceID)
}

func MatchRouter(Type, method, url string) *http_router.RouterInfo {
	mu.Lock()
	defer mu.Unlock()
	return r.Match(Type, method, url)
}
