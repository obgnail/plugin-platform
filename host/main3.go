package main

import (
	common "github.com/obgnail/plugin-platform/common/common_type"
	"plugin"
)

func GetPluginPath() string {
	a := `/go/src/github.com/obgnail/plugin-platform/plugin/example/backend/backend.so`
	return a
}

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

//func main2() {
//	localResource := local.New(nil)
//
//	type User struct {
//		UUID string `orm:"uuid"`
//		Name string `orm:"name"`
//	}
//	users := make([]*User, 0)
//
//	localResource.GetLocalDB().AsyncSelect("select * from user;", &users,
//		func(rawData []*common.RawData, descs []*common.ColumnDesc, pluginError common.PluginError, i interface{}) {
//			if pluginError != nil {
//				fmt.Println(pluginError)
//			}
//			err := localResource.GetLocalDB().Unmarshal(rawData, descs, i)
//			if err != nil {
//				fmt.Println(err)
//			}
//		})
//	fmt.Println(users)
//}
//
//func main3() {
//	p := GetPluginPath()
//	factory, err := getPluginFactory(p)
//	if err != nil {
//		panic(err)
//	}
//	plugin := factory()
//	localResource := local.New(plugin)
//	err = plugin.Assign(&mockInstanceDesc{}, localResource)
//	if err != nil {
//		panic(err)
//	}
//	res := plugin.Enable(common.LifeCycleRequest{})
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println("res", res)
//
//	fmt.Println("=================== end test ==================")
//}
