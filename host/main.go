package main

import (
	"fmt"
	common "github.com/obgnail/plugin-platform/common_type"
	"github.com/obgnail/plugin-platform/host/resource/local"
	"plugin"
)

func GetPluginPath() string {
	a := `/go/src/github.com/obgnail/plugin-platform/plugin/example/backend/backend.so`
	return a
}

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

func getPluginFactory(path string) (func() common.IPlugin, error) {
	p, err := plugin.Open(path)
	if err != nil {
		return nil, err
	}
	factory, err := p.Lookup("GetPlugin")
	if err != nil {
		return nil, err
	}

	return factory.(func() common.IPlugin), nil
}

func main() {
	localResource := local.New(nil)

	type User struct {
		UUID string `orm:"uuid"`
		Name string `orm:"name"`
	}
	users := make([]*User, 0)

	localResource.GetLocalDB().AsyncSelect("select * from user;", &users,
		func(rawData []*common.RawData, descs []*common.ColumnDesc, pluginError common.PluginError, i interface{}) {
			if pluginError != nil {
				fmt.Println(pluginError)
			}
			err := localResource.GetLocalDB().Unmarshal(rawData, descs, i)
			if err != nil {
				fmt.Println(err)
			}
		})
	fmt.Println(users)
}

func main3() {
	p := GetPluginPath()
	factory, err := getPluginFactory(p)
	if err != nil {
		panic(err)
	}
	plugin := factory()
	localResource := local.New(plugin)
	err = plugin.Assign(&mockInstanceDesc{}, localResource)
	if err != nil {
		panic(err)
	}
	res := plugin.Enable(common.LifeCycleRequest{})
	if err != nil {
		panic(err)
	}
	fmt.Println("res", res)

	fmt.Println("=================== end test ==================")
}
