package handler

import (
	"github.com/obgnail/plugin-platform/common/common_type"
	"sync"
)

type InstancePool struct {
	plugins sync.Map // map[instanceID]common_type.IPlugin
	running sync.Map // 在业务上host运行的插件 map[instanceID]common_type.IInstanceDescription
}

// ListInstances 列出所有的插件实例
func (pool *InstancePool) ListInstances() []common_type.IInstanceDescription {
	var resp = make([]common_type.IInstanceDescription, 0)
	pool.running.Range(func(k, v interface{}) bool {
		resp = append(resp, v.(common_type.IInstanceDescription))
		return true
	})
	return resp
}

func (pool *InstancePool) Add(instanceID string, plugin common_type.IPlugin) {
	pool.plugins.Store(instanceID, plugin)
}

func (pool *InstancePool) StartInstance(target common_type.IInstanceDescription) {
	instanceID := target.InstanceID()
	pool.running.Store(instanceID, target)
}

func (pool *InstancePool) StopInstance(instanceID string) {
	pool.running.Delete(instanceID)
}

func (pool *InstancePool) DeleteInstance(instanceID string) {
	pool.plugins.Delete(instanceID)
	pool.running.Delete(instanceID)
}

func (pool *InstancePool) GetInstance(instanceID string) bool {
	if _, ok := pool.running.Load(instanceID); ok {
		return true
	}
	return false
}

func (pool *InstancePool) GetPlugin(instanceID string) (plugin common_type.IPlugin, pluginDesc common_type.IInstanceDescription, exist bool) {
	val1, ok1 := pool.plugins.Load(instanceID)
	val2, ok2 := pool.running.Load(instanceID)

	if !(ok1 && ok2) {
		return nil, nil, false
	}

	return val1.(common_type.IPlugin), val2.(common_type.IInstanceDescription), true
}
