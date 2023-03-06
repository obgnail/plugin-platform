package handler

import (
	"github.com/obgnail/plugin-platform/common/common_type"
	"sync"
)

type InstancePool struct {
	mounted sync.Map // 已经挂载了的插件 map[instanceID]common_type.IPlugin
	running sync.Map // 在业务上host运行的插件 map[instanceID]common_type.IInstanceDescription
}

// ListRunning 列出所有运行的插件实例
func (pool *InstancePool) ListRunning() []common_type.IInstanceDescription {
	var resp []common_type.IInstanceDescription
	pool.running.Range(func(k, v interface{}) bool {
		resp = append(resp, v.(common_type.IInstanceDescription))
		return true
	})
	return resp
}

func (pool *InstancePool) ListMounted() []common_type.IPlugin {
	var resp []common_type.IPlugin
	pool.mounted.Range(func(k, v interface{}) bool {
		resp = append(resp, v.(common_type.IPlugin))
		return true
	})
	return resp
}

func (pool *InstancePool) AddMounted(instanceID string, plugin common_type.IPlugin) {
	pool.mounted.Store(instanceID, plugin)
}

func (pool *InstancePool) AddRunning(target common_type.IInstanceDescription) {
	instanceID := target.InstanceID()
	pool.running.Store(instanceID, target)
}

func (pool *InstancePool) DeleteRunning(instanceID string) {
	pool.running.Delete(instanceID)
}

func (pool *InstancePool) DeleteMountedAndRunning(instanceID string) {
	pool.mounted.Delete(instanceID)
	pool.running.Delete(instanceID)
}

func (pool *InstancePool) GetRunning(instanceID string) (common_type.IInstanceDescription, bool) {
	if val, ok := pool.running.Load(instanceID); ok {
		return val.(common_type.IInstanceDescription), true
	}
	return nil, false
}

func (pool *InstancePool) GetMounted(instanceID string) (common_type.IPlugin, bool) {
	if val, ok := pool.mounted.Load(instanceID); ok {
		return val.(common_type.IPlugin), true
	}
	return nil, false
}

func (pool *InstancePool) GetMountedAndRunning(instanceID string) (plugin common_type.IPlugin, pluginDesc common_type.IInstanceDescription, exist bool) {
	val1, ok1 := pool.mounted.Load(instanceID)
	val2, ok2 := pool.running.Load(instanceID)

	if !(ok1 && ok2) {
		return nil, nil, false
	}

	return val1.(common_type.IPlugin), val2.(common_type.IInstanceDescription), true
}
