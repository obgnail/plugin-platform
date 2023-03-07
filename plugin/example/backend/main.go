package main

import (
	"fmt"
	"github.com/obgnail/plugin-platform/common/common_type"
)

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

func (p *mockPlugin) Install(common_type.LifeCycleRequest) common_type.PluginError {
	event := p.resource.GetEventPublisher()
	cnd := []string{"project.task", "project.user"}
	er := event.Subscribe(cnd)
	if er != nil {
		panic(er)
	}
	er2 := event.Unsubscribe(cnd)
	if er2 != nil {
		panic(er2)
	}
	fmt.Println("-------------install-------------", cnd)
	return nil
}

func (p *mockPlugin) Enable(common_type.LifeCycleRequest) common_type.PluginError {
	//p.resource.GetLogger().Warn("this is warn )))(()()")

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
	fmt.Println("-------------enable------------- test localdb [select * from user;]")
	return nil
}

func (p *mockPlugin) Disable(common_type.LifeCycleRequest) common_type.PluginError {
	err := p.resource.GetWorkspace().CreateFile("mack_plugin_test.txt")
	if err != nil {
		panic(err)
	}
	fmt.Println("-------------disable-------------")
	return nil
}

func (p *mockPlugin) UnInstall(common_type.LifeCycleRequest) common_type.PluginError {
	fmt.Println("-------------uninstall-------------")
	return nil
}

func (p *mockPlugin) Upgrade(common_type.IVersion, common_type.LifeCycleRequest) common_type.PluginError {
	return nil
}

func (p *mockPlugin) OnEvent(eventType string, payload []byte) common_type.PluginError {
	fmt.Println("-------------onEvent-------------", eventType, string(payload))
	return common_type.NewPluginError(1, "unknown123XXXX")
}

func (p *mockPlugin) OnExternalHttpRequest(req *common_type.HttpRequest) *common_type.HttpResponse {
	fmt.Println("-------------OnExternalHttpRequest-------------", req.Url)
	body := "呵呵External"
	resp := &common_type.HttpResponse{
		Err:        nil,
		Headers:    req.Headers,
		Body:       []byte(body),
		StatusCode: 200,
	}
	return resp
}

func (p *mockPlugin) OnHttpCall(req *common_type.HttpRequest) (resp *common_type.HttpResponse) {
	fmt.Println("-------------OnHttpCall-------------", req.Url)
	body := "呵呵"
	resp = &common_type.HttpResponse{
		Err:        nil,
		Headers:    req.Headers,
		Body:       []byte(body),
		StatusCode: 200,
	}
	return resp
}

func (p *mockPlugin) OnConfigChange(configKey string, originValue, newValue []string) common_type.PluginError {
	fmt.Println("-------------OnConfigChange-------------", configKey, originValue[0], newValue[0])
	return common_type.NewPluginError(222, "unknown2222))))")
}

// TODO
func (p *mockPlugin) CheckState() common_type.PluginError {
	return nil
}

// TODO
func (p *mockPlugin) CheckCompatibility() common_type.PluginError {
	return nil
}
