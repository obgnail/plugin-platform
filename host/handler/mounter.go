package handler

import (
	"fmt"
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/errors"
	"github.com/obgnail/plugin-platform/host/resource/common"
	"github.com/obgnail/plugin-platform/host/resource/local"
	"github.com/obgnail/plugin-platform/host/resource/release"
	"github.com/obgnail/plugin-platform/platform/service/utils"
	"plugin"
)

type PluginMounter struct {
	sender  common.Sender
	isLocal bool
}

func NewMounter(sender common.Sender, isLocal bool) *PluginMounter {
	return &PluginMounter{isLocal: isLocal, sender: sender}
}

func (m *PluginMounter) Setup(unset common_type.IPlugin, description common_type.IInstanceDescription) (
	setup common_type.IPlugin, err common_type.PluginError) {

	if unset == nil {
		desc := description.PluginDescription()
		var er error
		unset, er = m.CreatePlugin(desc.ApplicationID(), desc.ApplicationVersion().VersionString())
		if er != nil {
			return nil, common_type.NewPluginError(common_type.GetInstanceFailure, er.Error())
		}
	}

	resources := m.GetResources(unset)

	setup, err = common_type.SetupPlugin(unset, description, resources)
	if err != nil {
		return nil, err
	}

	return setup, nil
}

func (m *PluginMounter) CreatePlugin(appID, appVersion string) (common_type.IPlugin, error) {
	path := utils.GetPluginSoFile(appID, appVersion)
	_plugin, err := getPlugin(path)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return _plugin, nil
}

func (m *PluginMounter) GetResources(plugin common_type.IPlugin) common_type.IResources {
	if m.isLocal {
		return &local.Resource{Plugin: plugin}
	}
	return &release.Resource{Plugin: plugin, Sender: m.sender}
}

func getPlugin(path string) (common_type.IPlugin, error) {
	p, err := plugin.Open(path)
	if err != nil {
		return nil, err
	}
	factory, err := p.Lookup("GetPlugin")
	if err != nil {
		return nil, err
	}
	f, ok := factory.(func() common_type.IPlugin)
	if !ok {
		return nil, fmt.Errorf("not implement GetPlugin() function")
	}
	return f(), nil
}
