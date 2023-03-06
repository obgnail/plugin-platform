package router

import (
	"github.com/obgnail/plugin-platform/common/errors"
	"github.com/obgnail/plugin-platform/platform/service/common"
	"sync"
)

type router struct {
	method string
	path   string
}

type routerInfo struct {
	functionName string
	instanceUUID string
}

type PluginRouter struct {
	mu        sync.RWMutex           // protect below
	apiCenter map[router]struct{}    // 全局路由列表: 只关注路由本身;使用map是为了快速判断冲突
	routerMap map[router]*routerInfo // 全局路由列表: 维护路由和插件实例的关系
}

func NewRouter() *PluginRouter {
	return &PluginRouter{apiCenter: make(map[router]struct{}), routerMap: make(map[router]*routerInfo)}
}

func (r *PluginRouter) Add(apis []*common.Api, instanceID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	routers, routerInfos := r.zip(apis, instanceID)
	if err := r.checkConflict(routers, routerInfos); err != nil {
		return errors.Trace(err)
	}
	r.saveApi(routers, routerInfos)
	return nil
}

func (r *PluginRouter) zip(apis []*common.Api, instanceID string) (routers []router, routerInfos []*routerInfo) {
	for _, api := range apis {
		for _, method := range api.Methods {
			routers = append(routers, router{method: method, path: api.Url})
			routerInfos = append(routerInfos, &routerInfo{functionName: api.Function, instanceUUID: instanceID})
		}
	}
	return
}

// TODO check conflict
func (r *PluginRouter) checkConflict(routers []router, routerInfos []*routerInfo) error {
	return nil
}

func (r *PluginRouter) saveApi(routers []router, routerInfos []*routerInfo) {
	for index, _r := range routers {
		r.apiCenter[_r] = struct{}{}
		r.routerMap[_r] = routerInfos[index]
	}
}
