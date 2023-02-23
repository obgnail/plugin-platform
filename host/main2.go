package main

import (
	"fmt"
	common "github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/log"
	"github.com/obgnail/plugin-platform/host/resource/release"
)

func main() {
	plugin := &mockPlugin{}
	resource := release.NewReleaseResource(plugin)
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
func (i *mockInstanceDesc) PluginDescription() common.PluginDescriptor {
	return &mockPluginDescriptor{}
}

type mockPluginDescriptor struct{}

func (i *mockPluginDescriptor) ApplicationID() string            { return "ApplicationID123" }
func (i *mockPluginDescriptor) Name() string                     { return "ApplicationID123" }
func (i *mockPluginDescriptor) Language() string                 { return "Language123" }
func (i *mockPluginDescriptor) LanguageVersion() common.IVersion { return common.NewVersion(1, 2, 3) }
func (i *mockPluginDescriptor) ApplicationVersion() common.IVersion {
	return common.NewVersion(1, 2, 3)
}
func (i *mockPluginDescriptor) HostVersion() common.IVersion      { return common.NewVersion(1, 2, 3) }
func (i *mockPluginDescriptor) MinSystemVersion() common.IVersion { return common.NewVersion(1, 2, 3) }

////////////////////////////////////////////////////////////////////////////////////////////////
var _ common.IPlugin = (*mockPlugin)(nil)

type mockPlugin struct {
	descriptor common.IInstanceDescription
	resource   common.IResources
}

func GetPlugin() common.IPlugin {
	return &mockPlugin{}
}
func (p *mockPlugin) GetPluginDescription() common.IInstanceDescription {
	return p.descriptor
}

func (p *mockPlugin) Assign(pd common.IInstanceDescription, resource common.IResources) common.PluginError {
	p.descriptor = pd
	p.resource = resource
	return nil
}

func (p *mockPlugin) Enable(common.LifeCycleRequest) common.PluginError {
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

func (p *mockPlugin) Disable(common.LifeCycleRequest) common.PluginError {
	err := p.resource.GetWorkspace().CreateFile("mack_plugin_test.txt")
	if err != nil {
		log.ErrorDetails(err)
	}
	return nil
}

func (p *mockPlugin) Start(common.LifeCycleRequest) common.PluginError {
	err := p.resource.GetWorkspace().WriteBytes("mack_plugin_test.txt", []byte("1234567654321"))
	if err != nil {
		log.ErrorDetails(err)
	}
	return nil
}

func (p *mockPlugin) Stop(common.LifeCycleRequest) common.PluginError {
	return nil
}

// TODO
func (p *mockPlugin) CheckState() common.PluginError {
	return nil
}

// TODO
func (p *mockPlugin) CheckCompatibility() common.PluginError {
	return nil
}

func (p *mockPlugin) Install(common.LifeCycleRequest) common.PluginError {
	return nil
}

func (p *mockPlugin) UnInstall(common.LifeCycleRequest) common.PluginError {
	return nil
}

func (p *mockPlugin) Upgrade(common.IVersion, common.LifeCycleRequest) common.PluginError {
	return nil
}

func (p *mockPlugin) OnEvent(eventType string, payload interface{}) common.PluginError {
	return nil
}

func (p *mockPlugin) OnExternalHttpRequest(request *common.HttpRequest) *common.HttpResponse {
	return nil
}
