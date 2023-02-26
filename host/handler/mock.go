package handler

import (
	"github.com/obgnail/plugin-platform/common/common_type"
)

var _ common_type.IPlugin = (*MockPlugin)(nil)
var _ common_type.IInstanceDescription = (*MockInstanceDesc)(nil)
var _ common_type.IPluginDescriptor = (*MockPluginDescriptor)(nil)

type MockInstanceDesc struct {
	instanceID       string
	pluginDescriptor *MockPluginDescriptor
}

func (i *MockInstanceDesc) InstanceID() string { return i.instanceID }
func (i *MockInstanceDesc) PluginDescription() common_type.IPluginDescriptor {
	return i.pluginDescriptor
}

type MockPluginDescriptor struct {
	appID     string
	name      string
	lang      string
	langVer   string
	appVer    string
	hostVer   string
	minSysVer string
}

func (i *MockPluginDescriptor) ApplicationID() string { return i.appID }
func (i *MockPluginDescriptor) Name() string          { return i.name }
func (i *MockPluginDescriptor) Language() string      { return i.lang }
func (i *MockPluginDescriptor) LanguageVersion() common_type.IVersion {
	result, _ := common_type.ParseVersionString(i.langVer)
	return result
}
func (i *MockPluginDescriptor) ApplicationVersion() common_type.IVersion {
	result, _ := common_type.ParseVersionString(i.appVer)
	return result
}
func (i *MockPluginDescriptor) HostVersion() common_type.IVersion {
	result, _ := common_type.ParseVersionString(i.hostVer)
	return result
}
func (i *MockPluginDescriptor) MinSystemVersion() common_type.IVersion {
	result, _ := common_type.ParseVersionString(i.minSysVer)
	return result
}

type MockPlugin struct {
	common_type.IPlugin
}

func (w *MockPlugin) Setup(p common_type.IPlugin, desc common_type.IInstanceDescription, res common_type.IResources) common_type.PluginError {
	w.IPlugin = p
	err := w.IPlugin.Assign(desc, res)
	return err
}

func SetupPlugin(plugin common_type.IPlugin, desc common_type.IInstanceDescription, res common_type.IResources) (common_type.IPlugin, common_type.PluginError) {
	mock := &MockPlugin{}
	err := mock.Setup(plugin, desc, res)
	return mock, err
}
