package platform

import (
	"encoding/json"
	"github.com/obgnail/plugin-platform/common/errors"
	"github.com/obgnail/plugin-platform/platform/conn/hub/router"
	"github.com/obgnail/plugin-platform/platform/conn/hub/router/http_router"
	"github.com/obgnail/plugin-platform/platform/service/common"
	"sync"
)

var mu sync.Mutex
var route *router.PluginRouter

type Plugin struct {
	UUID        string    `json:"uuid"`
	Name        string    `json:"name"`
	Version     string    `json:"version"`
	LifeStage   int       `json:"life_stage"`
	Description string    `json:"description"`
	Routers     []*Router `json:"routers"`
}

type Router struct {
	Type   string   `json:"type"`
	Method []string `json:"methods"`
	Url    string   `json:"url"`
}

func UnmarshalRouterList(resp []byte) ([]*Plugin, error) {
	s := struct {
		Data []*Plugin `json:"data"`
	}{}
	if err := json.Unmarshal(resp, &s); err != nil {
		return nil, errors.Trace(err)
	}

	return s.Data, nil
}

func Register(plugins []*Plugin) error {
	mu.Lock()
	defer mu.Unlock()

	route = router.NewRouter()

	var apis []*common.Api

	for _, plugin := range plugins {
		if plugin.LifeStage != common.PluginStatusRunning {
			continue
		}
		for _, r := range plugin.Routers {
			apis = append(apis, &common.Api{
				Type:    r.Type,
				Methods: r.Method,
				Url:     r.Url,
			})
		}
		if err := route.Add(apis, plugin.UUID); err != nil {
			return errors.Trace(err)
		}
	}
	return nil
}

func MatchRouter(Type, method, url string) *http_router.RouterInfo {
	mu.Lock()
	defer mu.Unlock()
	return route.Match(Type, method, url)
}
