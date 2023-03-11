package platform

import (
	"encoding/json"
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/errors"
	"github.com/obgnail/plugin-platform/common/log"
	"github.com/obgnail/plugin-platform/platform/conn/hub/ability"
	"github.com/obgnail/plugin-platform/platform/conn/hub/router"
	"github.com/obgnail/plugin-platform/platform/conn/hub/router/http_router"
	"github.com/obgnail/plugin-platform/platform/service/common"
	"sync"
)

var (
	mu         sync.Mutex // protect below
	routeHub   *router.PluginRouter
	abilityHub *ability.Ability
)

func registerHub(plugins []*Plugin) error {
	mu.Lock()
	defer mu.Unlock()

	routeHub = router.NewRouter()
	abilityHub = ability.NewAbility()

	err := rangeRunningPlugin(plugins, func(plugin *Plugin) error {
		if err := routeHub.Register(plugin.UUID, plugin.Routers); err != nil {
			return errors.Trace(err)
		}
		abilityHub.Register(plugin.UUID, plugin.Abilities)
		return nil
	})

	if err != nil {
		return errors.Trace(err)
	}

	return nil
}

func MatchRouter(Type, method, url string) *http_router.RouterInfo {
	mu.Lock()
	defer mu.Unlock()
	if routeHub == nil {
		log.Error("router not init")
		return nil
	}
	return routeHub.Match(Type, method, url)
}

func ExecuteAbility(instanceID, abilityID, abilityType, abilityFuncKey string, arg []byte) (chan *common_type.AbilityResponse, error) {
	return abilityHub.Execute(instanceID, abilityID, abilityType, abilityFuncKey, arg)
}

type Plugin struct {
	UUID        string            `json:"uuid"`
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	LifeStage   int               `json:"life_stage"`
	Description string            `json:"description"`
	Routers     []*common.Api     `json:"routers"`
	Abilities   []*common.Ability `json:"abilities"`
}

func rangeRunningPlugin(plugins []*Plugin, f func(plugin *Plugin) error) error {
	for _, plugin := range plugins {
		if plugin.LifeStage != common.PluginStatusRunning {
			continue
		}
		if err := f(plugin); err != nil {
			return errors.Trace(err)
		}
	}
	return nil
}

func unmarshalPluginList(resp []byte) ([]*Plugin, error) {
	s := struct {
		Data []*Plugin `json:"data"`
	}{}
	if err := json.Unmarshal(resp, &s); err != nil {
		return nil, errors.Trace(err)
	}

	return s.Data, nil
}
