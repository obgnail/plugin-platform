package router

import (
	"github.com/obgnail/plugin-platform/common/errors"
	"github.com/obgnail/plugin-platform/platform/conn/hub/router/http_router"
	"github.com/obgnail/plugin-platform/platform/service/common"
	"strings"
)

type PluginRouter struct {
	route *http_router.Router
}

func NewRouter() *PluginRouter {
	return &PluginRouter{route: http_router.New()}
}

func (r *PluginRouter) Add(apis []*common.Api, instanceID string) error {
	backup := r.route.DeepCopy()

	for _, api := range apis {
		for _, method := range api.Methods {
			router := &http_router.RouterInfo{FunctionName: api.Function, InstanceUUID: instanceID}

			if err := r.route.AddRoute(strings.ToLower(api.Type), strings.ToLower(method), api.Url, router); err != nil {
				r.route = backup // 失败时还原
				return errors.Trace(err)
			}
		}
	}
	return nil
}

func (r *PluginRouter) Delete(instanceID string) {
	// 将instanceID对应的所有handle删掉,即是删除
	r.route.RangeRoute(func(Type, method string, route *http_router.RouterInfo) {
		if route != nil && route.InstanceUUID == instanceID {
			route = nil
		}
	})
}

func (r *PluginRouter) Match(Type, method, url string) *http_router.RouterInfo {
	return r.route.GetRouter(strings.ToLower(Type), strings.ToLower(method), url)
}
