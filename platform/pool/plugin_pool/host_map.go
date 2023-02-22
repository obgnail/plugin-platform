package plugin_pool

import (
	"fmt"
	"sync"
)

// HostBasicInformation 每个host支持的插件列表
type HostBasicInformation struct {
	//PluginDescriptors []*protocol.PluginDescriptor
	//HostDescriptor    *protocol.HostDescriptor
}

type HostMap struct {
	m sync.Map
}

func NewHostMap() *HostMap {
	return new(HostMap)
}

func (Map *HostMap) New(hostID string, hostBasicInformation *HostBasicInformation) {
	Map.m.Store(hostID, hostBasicInformation)
}

func (Map *HostMap) One(hostID string) (*HostBasicInformation, error) {
	v, ok := Map.m.Load(hostID)
	if !ok {
		return nil, fmt.Errorf("no such host: %s", hostID)
	}
	p := v.(*HostBasicInformation)
	return p, nil
}

func (Map *HostMap) All() []*HostBasicInformation {
	var resp = make([]*HostBasicInformation, 0)
	Map.m.Range(func(k, v interface{}) bool {
		p := v.(*HostBasicInformation)
		resp = append(resp, p)
		return true
	})
	return resp
}

func (Map *HostMap) Delete(hostID string) {
	Map.m.Delete(hostID)
}

//插件生命周期使用 获取匹配的hostIDs
//func (Map *HostMap) GetHostIDs(language, languageVersion, hostVerison string) []*HostInfo {
//re := make([]*HostInfo, 0)
//hostMaps := Map.All()
//for _, host := range hostMaps {
//	hostLanguage := host.HostDescriptor.GetLanguage()
//	hostLanguageVersion := connection.GetVersionString(host.HostDescriptor.GetLanguageVersion())
//	HostVersion := connection.GetVersionString(host.HostDescriptor.GetHostVersion())
//	if language == hostLanguage && hostLanguageVersion == languageVersion && HostVersion == hostVerison {
//		tmp := &HostInfo{}
//		tmp.HostID = host.HostDescriptor.GetHostID()
//		tmp.HostVersion = host.HostDescriptor.GetHostVersion().String()
//		tmp.HostName = host.HostDescriptor.GetName()
//		tmp.HostDescriptor = host.HostDescriptor
//		re = append(re, tmp)
//	}
//}
//	return re
//}
