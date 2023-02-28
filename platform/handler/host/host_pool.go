package host

import (
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/message_utils"
	"github.com/obgnail/plugin-platform/common/protocol"
	"sync"
)

type Pool struct {
	alive sync.Map // 在业务上运行的host map[hostID]common_type.IHost
}

func NewPool() *Pool {
	return &Pool{}
}

func (pool *Pool) Add(host common_type.IHost) {
	id := host.GetInfo().ID
	pool.alive.Store(id, host)
}

func (pool *Pool) Delete(host common_type.IHost) {
	id := host.GetInfo().ID
	pool.alive.Delete(id)
}

func (pool *Pool) Exist(hostID string) bool {
	if _, ok := pool.alive.Load(hostID); ok {
		return true
	}
	return false
}

func (pool *Pool) GetAll() []common_type.IHost {
	var result []common_type.IHost
	pool.alive.Range(func(key, value any) bool {
		result = append(result, value.(common_type.IHost))
		return true
	})
	return result
}

func (pool *Pool) Range(f func(hostID string, host common_type.IHost) bool) {
	pool.alive.Range(func(key, value any) bool {
		return f(key.(string), value.(common_type.IHost))
	})
}

func NewHost(msg *protocol.ControlMessage_HostReportMessage, status common_type.HostStatus) *common_type.MockHost {
	_host := msg.GetHost()
	plugins := msg.GetInstanceList()

	_plugins := make(map[string]common_type.IInstanceDescription, len(plugins))
	for _, info := range plugins {
		instanceID := info.GetInstanceID()
		appDesc := info.GetApplication()

		_plugins[instanceID] = &common_type.MockInstanceDesc{
			PluginInstanceID: instanceID,
			PluginDescriptor: &common_type.MockPluginDescriptor{
				AppID:      appDesc.ApplicationID,
				PluginName: appDesc.Name,
				Lang:       appDesc.Language,
				LangVer:    message_utils.VersionPb2String(appDesc.LanguageVersion),
				AppVer:     message_utils.VersionPb2String(appDesc.ApplicationVersion),
				HostVer:    message_utils.VersionPb2String(appDesc.HostVersion),
				MinSysVer:  message_utils.VersionPb2String(appDesc.MinSystemVersion),
			},
		}
	}

	h := &common_type.MockHost{
		Info: common_type.HostInfo{
			ID:               _host.GetHostID(),
			Name:             _host.GetName(),
			Version:          message_utils.VersionPb2String(_host.GetHostVersion()),
			MinSystemVersion: message_utils.VersionPb2String(_host.GetMinSystemVersion()),
			Language:         _host.GetLanguage(),
			LanguageVersion:  message_utils.VersionPb2String(_host.GetLanguageVersion()),
			Plugins:          _plugins,
		},
		Status: status,
	}
	return h
}
