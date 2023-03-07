package handler

import (
	"fmt"
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/errors"
	"github.com/obgnail/plugin-platform/host/resource/common"
	"github.com/obgnail/plugin-platform/platform/service/utils"
	"plugin"
)

func MountPlugin(
	iPlugin common_type.IPlugin,
	iDescription common_type.IInstanceDescription,
	resourceGetter common.ResourceFactor,
	sender common.Sender,
) (common_type.IPlugin, common_type.PluginError) {

	// 如果不为nil则复用,否则创建一个新的
	if iPlugin == nil {
		var err error
		desc := iDescription.PluginDescription()
		iPlugin, err = createPlugin(desc.ApplicationID(), desc.ApplicationVersion().VersionString())
		if err != nil {
			return nil, common_type.NewPluginError(common_type.GetInstanceFailure, err.Error())
		}
	}

	iResource := resourceGetter.GetResource(iPlugin, sender)
	if err := iPlugin.Assign(iDescription, iResource); err != nil {
		return nil, err
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
