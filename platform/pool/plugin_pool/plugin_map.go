package plugin_pool

import (
	"fmt"
	"sync"
)

type PluginMap struct {
	m sync.Map
}

func NewPluginMap() *PluginMap {
	return new(PluginMap)
}

func (Map *PluginMap) New(pluginUUID string, plugin PluginInterface) {
	Map.m.Store(pluginUUID, plugin)
}

func (Map *PluginMap) One(pluginUUID string) (PluginInterface, error) {
	v, ok := Map.m.Load(pluginUUID)
	if !ok {
		return nil, fmt.Errorf("no such plugin: %s", pluginUUID)
	}
	p := v.(PluginInterface)
	return p, nil
}

func (Map *PluginMap) All() []PluginInterface {
	var resp = make([]PluginInterface, 0)
	Map.m.Range(func(k, v interface{}) bool {
		p := v.(PluginInterface)
		resp = append(resp, p)
		return true
	})
	return resp
}

func (Map *PluginMap) Delete(pluginUUID string) {
	Map.m.Delete(pluginUUID)
}
