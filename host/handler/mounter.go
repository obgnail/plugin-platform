package handler

import (
	"fmt"
	"github.com/obgnail/plugin-platform/common/common_type"
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

func (m *PluginMounter) Mount(Plugin common_type.IPlugin, instanceDesc common_type.IInstanceDescription) (common_type.IPlugin, common_type.PluginError) {
	var err common_type.PluginError

	if Plugin == nil {
		desc := instanceDesc.PluginDescription()
		Plugin, err = m.CreatePlugin(desc.ApplicationID(), desc.ApplicationVersion().VersionString())
		if err != nil {
			return nil, err
		}
	}

	resources := m.GetResources(Plugin)

	setupPlugin, err := common_type.SetupPlugin(Plugin, instanceDesc, resources)
	if err != nil {
		return nil, err
	}

	return setupPlugin, nil
}

func (m *PluginMounter) CreatePlugin(appID, appVersion string) (common_type.IPlugin, common_type.PluginError) {
	path := utils.GetPluginSoFile(appID, appVersion)
	_plugin, err := getPlugin(path)
	if err != nil {
		e := common_type.NewPluginError(common_type.GetInstanceFailure, fmt.Sprintf("appid: %s plugin so file not found: %v", appID, err))
		return nil, e
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
