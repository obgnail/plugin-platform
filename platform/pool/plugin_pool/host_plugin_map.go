package plugin_pool

import (
	"fmt"
	"sync"
)

// PluginInstanceInfo 每个插件实例的一些信息
type PluginInstanceInfo struct {
	InstanceID string
	IsStop     bool //用来判断是否用户操作的stop
}

type HostPluginMap struct {
	m sync.Map
}

// HostPluginInstance 每个host上的插件实例
type HostPluginInstance struct {
	//HostPluginInstanceList []*protocol.PluginInstanceDescriptor
	//HostPluginInstanceInfo []*PluginInstanceInfo
	//HostDescriptor         *protocol.HostDescriptor
}

func NewHostPluginMap() *HostPluginMap {
	return new(HostPluginMap)
}

func (Map *HostPluginMap) New(hostID string, hostPluginInstance *HostPluginInstance) {
	Map.m.Store(hostID, hostPluginInstance)
}

func (Map *HostPluginMap) One(hostID string) (*HostPluginInstance, error) {
	v, ok := Map.m.Load(hostID)
	if !ok {
		return nil, fmt.Errorf("no such host: %s", hostID)
	}
	p := v.(*HostPluginInstance)
	return p, nil
}

func (Map *HostPluginMap) All() []*HostPluginInstance {
	var resp = make([]*HostPluginInstance, 0)
	Map.m.Range(func(k, v interface{}) bool {
		p := v.(*HostPluginInstance)
		resp = append(resp, p)
		return true
	})
	return resp
}

func (Map *HostPluginMap) Delete(hostID string) {
	Map.m.Delete(hostID)
}

//
////重启机制 通过插件实例ID获取hosts
//func (Map *HostPluginMap) GetHostIDsByPluginInstanceID(instanceID string) []*HostInfo {
//	re := make([]*HostInfo,0)
//	hostPluginInstanceMaps := Map.All()
//	for _,hostPluginInstance := range hostPluginInstanceMaps {
//		for _,pluginInstances := range hostPluginInstance.HostPluginInstanceList{
//			if pluginInstances.GetInstanceID() == instanceID{
//				tmp := &HostInfo{}
//				tmp.HostID         = pluginInstances.GetHostID()
//				tmp.HostName       = hostPluginInstance.HostDescriptor.GetName()
//				tmp.HostVersion    = hostPluginInstance.HostDescriptor.GetHostVersion().String()
//				tmp.HostDescriptor = hostPluginInstance.HostDescriptor
//				re = append(re,tmp)
//			}
//		}
//	}
//	return re
//}
//
////重启机制 通过hostID获取插件实例ID 排除掉stop的插件
//func (Map *HostPluginMap) GetPluginInstanceIDsByHostID(hostID string) ([]*protocol.PluginInstanceDescriptor,error) {
//	re := make([]*protocol.PluginInstanceDescriptor,0)
//	hostPluginInstance,err := Map.One(hostID)
//	if err != nil {
//		return re,err
//	}
//	var pluginInstanceStopInfo = make(map[string]string)
//	for _,v := range hostPluginInstance.HostPluginInstanceInfo{
//		if v.IsStop == true{
//			pluginInstanceStopInfo[v.InstanceID] = v.InstanceID
//		}
//	}
//	for _,pluginInstance := range hostPluginInstance.HostPluginInstanceList {
//		if _,ok:=pluginInstanceStopInfo[pluginInstance.InstanceID];!ok{
//			re = append(re,pluginInstance)
//		}
//	}
//	return re,nil
//}
//
////通过插件实例ID 删除
//func (Map *HostPluginMap) DeletePluginInstanceByInstanceIDAndHostID(instanceID,hostID string) error {
//	hostPluginInstance,err := Map.One(hostID)
//	if err != nil {
//		return err
//	}
//	tmp := &HostPluginInstance{
//		HostPluginInstanceList:make([]*protocol.PluginInstanceDescriptor,0),
//		HostPluginInstanceInfo:make([]*PluginInstanceInfo,0),
//	}
//	tmp.HostDescriptor = hostPluginInstance.HostDescriptor
//	for index,pluginInstance := range hostPluginInstance.HostPluginInstanceList {
//		if pluginInstance.GetInstanceID() == instanceID {
//			hostPluginInstance.HostPluginInstanceList = append(hostPluginInstance.HostPluginInstanceList[0:index],hostPluginInstance.HostPluginInstanceList[index+1:]...)
//			tmp.HostPluginInstanceList = hostPluginInstance.HostPluginInstanceList
//		}
//	}
//
//	for index,pluginInstance := range hostPluginInstance.HostPluginInstanceInfo {
//		if pluginInstance.InstanceID == instanceID {
//			hostPluginInstance.HostPluginInstanceInfo = append(hostPluginInstance.HostPluginInstanceInfo[0:index],hostPluginInstance.HostPluginInstanceInfo[index+1:]...)
//			tmp.HostPluginInstanceInfo = hostPluginInstance.HostPluginInstanceInfo
//		}
//	}
//	Map.New(hostID,tmp)
//	return nil
//}
//
////修改某个host下的实例ID
//func (Map *HostPluginMap) UpdatePluginInstanceStopByInstanceIDAndHostID(instanceID,hostID string) error {
//	hostPluginInstance,err := Map.One(hostID)
//	if err != nil {
//		return err
//	}
//	for _,pluginInstance := range hostPluginInstance.HostPluginInstanceInfo {
//		if pluginInstance.InstanceID == instanceID {
//			pluginInstance.IsStop = true
//		}
//	}
//	Map.New(hostID,hostPluginInstance)
//	return nil
//}
//
////获取某个host下的pluginInstanceList
//func (Map *HostPluginMap) GetHostPluginInstanceList(hostID string) (map[string]*PluginInstanceInfo,error) {
//	hostPluginInstance,err := Map.One(hostID)
//	if err != nil {
//		return nil,err
//	}
//	var retMap = make(map[string]*PluginInstanceInfo,0)
//	for _,pluginInstance := range hostPluginInstance.HostPluginInstanceInfo {
//		retMap[pluginInstance.InstanceID] = pluginInstance
//	}
//	return retMap,nil
//}
