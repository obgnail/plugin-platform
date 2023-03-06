package router

import "github.com/obgnail/plugin-platform/platform/service/common"

var r *PluginRouter

func InitRouter() error {
	r = NewRouter()
	return nil
}

func RegisterRouter(apis []*common.Api, instanceID string) error {
	return r.Add(apis, instanceID)
}
