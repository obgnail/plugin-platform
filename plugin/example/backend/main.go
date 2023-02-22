package main

import (
	"fmt"
	common "github.com/obgnail/plugin-platform/common_type"
	"github.com/obgnail/plugin-platform/utils/log"
)

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
