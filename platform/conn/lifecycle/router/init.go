package router

import "github.com/obgnail/plugin-platform/platform/service/common"

var r *PluginRouter

func InitRouter() {
	r = NewRouter()
}

func RegisterRouter(apis []*common.Api, instanceID string) error {
	return nil
	//return r.Add(apis, instanceID)
}
