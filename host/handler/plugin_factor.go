package handler

import (
	"fmt"
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/errors"
	"github.com/obgnail/plugin-platform/platform/service/utils"
	"plugin"
)

func CreatePlugin(iDescription common_type.IInstanceDescription) (common_type.IPlugin, common_type.PluginError) {
	desc := iDescription.PluginDescription()
	iPlugin, err := createPlugin(desc.ApplicationID(), desc.ApplicationVersion().VersionString())
	if err != nil {
		return nil, common_type.NewPluginError(common_type.GetInstanceFailure, err.Error())
	}
	return iPlugin, nil
}

func createPlugin(appID, appVersion string) (common_type.IPlugin, error) {
	path := utils.GetPluginSoFile(appID, appVersion)
	_plugin, err := getPlugin(path)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return _plugin, nil
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
