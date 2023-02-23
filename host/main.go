package main

import (
	"fmt"
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/log"
	"github.com/obgnail/plugin-platform/host/resource/release"
)

func main() {
	plugin := &mockPlugin{}
	resource := release.NewResource(plugin)
	instanceDesc := &mockInstanceDesc{}

	log.Info("OK")
	err := plugin.Assign(instanceDesc, resource)
	if err != nil {
		panic(err)
	}
	workSpace := resource.GetWorkspace()
	e := workSpace.CreateFile("qwe_____.txt")
	if e != nil {
		panic(e)
	}
}

////////////////////////////
type mockInstanceDesc struct{}

func (i *mockInstanceDesc) InstanceID() string { return "InstanceID123" }
func (i *mockInstanceDesc) PluginDescription() common_type.PluginDescriptor {
	return &mockPluginDescriptor{}
}

type mockPluginDescriptor struct{}

func (i *mockPluginDescriptor) ApplicationID() string { return "ApplicationID123" }
func (i *mockPluginDescriptor) Name() string          { return "ApplicationID123" }
func (i *mockPluginDescriptor) Language() string      { return "Language123" }
func (i *mockPluginDescriptor) LanguageVersion() common_type.IVersion {
	return common_type.NewVersion(1, 2, 3)
}
func (i *mockPluginDescriptor) ApplicationVersion() common_type.IVersion {
	return common_type.NewVersion(1, 2, 3)
}
func (i *mockPluginDescriptor) HostVersion() common_type.IVersion {
	return common_type.NewVersion(1, 2, 3)
}
func (i *mockPluginDescriptor) MinSystemVersion() common_type.IVersion {
	return common_type.NewVersion(1, 2, 3)
}

////////////////////////////////////////////////////////////////////////////////////////////////
var _ common_type.IPlugin = (*mockPlugin)(nil)

type mockPlugin struct {
	descriptor common_type.IInstanceDescription
	resource   common_type.IResources
}

func GetPlugin() common_type.IPlugin {
	return &mockPlugin{}
}
func (p *mockPlugin) GetPluginDescription() common_type.IInstanceDescription {
	return p.descriptor
}

func (p *mockPlugin) Assign(pd common_type.IInstanceDescription, resource common_type.IResources) common_type.PluginError {
	p.descriptor = pd
	p.resource = resource
	return nil
}

func (p *mockPlugin) Enable(common_type.LifeCycleRequest) common_type.PluginError {
	fmt.Println("enable-------------")
	//err := p.resource.GetLocalDB().ImportSQL("/go/src/github.com/obgnail/plugin-platform/plugin/example/config/plugin.sql")
	//if err != nil {
	//	fmt.Println(err)
	//}
	rawData, colDesc, err := p.resource.GetLocalDB().Select("select * from user;")
	if err != nil {
		fmt.Println(err)
	}

	type User struct {
		uuid string
		name string
	}
	var users []*User

	err = p.resource.GetLocalDB().Unmarshal(rawData, colDesc, users)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(err)
	return nil
}

func (p *mockPlugin) Disable(common_type.LifeCycleRequest) common_type.PluginError {
	err := p.resource.GetWorkspace().CreateFile("mack_plugin_test.txt")
	if err != nil {
		log.ErrorDetails(err)
	}
	return nil
}

func (p *mockPlugin) Start(common_type.LifeCycleRequest) common_type.PluginError {
	err := p.resource.GetWorkspace().WriteBytes("mack_plugin_test.txt", []byte("1234567654321"))
	if err != nil {
		log.ErrorDetails(err)
	}
	return nil
}

func (p *mockPlugin) Stop(common_type.LifeCycleRequest) common_type.PluginError {
	return nil
}

// TODO
func (p *mockPlugin) CheckState() common_type.PluginError {
	return nil
}

// TODO
func (p *mockPlugin) CheckCompatibility() common_type.PluginError {
	return nil
}

func (p *mockPlugin) Install(common_type.LifeCycleRequest) common_type.PluginError {
	return nil
}

func (p *mockPlugin) UnInstall(common_type.LifeCycleRequest) common_type.PluginError {
	return nil
}

func (p *mockPlugin) Upgrade(common_type.IVersion, common_type.LifeCycleRequest) common_type.PluginError {
	return nil
}

func (p *mockPlugin) OnEvent(eventType string, payload interface{}) common_type.PluginError {
	return nil
}

func (p *mockPlugin) OnExternalHttpRequest(request *common_type.HttpRequest) *common_type.HttpResponse {
	return nil
}
