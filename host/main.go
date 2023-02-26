package main

import (
	"encoding/json"
	"fmt"
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/log"
	"github.com/obgnail/plugin-platform/host/resource/release"
	"time"
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
	//testWorkSpace(plugin)
	//testDB(plugin)
	//testNetwork(plugin)
	//testLog(plugin)
	testEvent(plugin)
}

func testEvent(plugin *mockPlugin) {
	event := plugin.resource.GetEventPublisher()
	cnd := []string{"project.task", "project.user"}
	er := event.Subscribe(cnd)
	if er != nil {
		panic(er)
	}
	er2 := event.Unsubscribe(cnd)
	if er2 != nil {
		panic(er2)
	}
}

func testLog(plugin *mockPlugin) {
	logger := plugin.resource.GetLogger()
	logger.Trace("****** trace ******")
	logger.Info("****** info ******")
	logger.Warn("****** warn ******")
	logger.Error("****** error ******")
}

func testNetwork(plugin *mockPlugin) {
	sysDB := plugin.resource.GetSysDB()
	rawData, cols, err := sysDB.Select("plugin_platform", "select * from plugin_package;")
	if err != nil {
		panic(err)
	}

	type Package struct {
		ID         string `orm:"id"`
		AppUUID    string `orm:"app_uuid"`
		Name       string `orm:"name"`
		Size       int    `orm:"size"`
		Version    string `orm:"version"`
		CreateTime int    `orm:"create_time"`
		UpdateTime int    `orm:"update_time"`
		Deleted    int    `orm:"deleted"`
	}
	ps := make([]*Package, 0)
	e := sysDB.Unmarshal(rawData, cols, &ps)
	if e != nil {
		panic(e)
	}

	result, er := json.Marshal(ps)
	if er != nil {
		panic(er)
	}

	api := plugin.resource.GetAPICore()
	req := &common_type.HttpRequest{
		Method:   "post",
		QueryMap: nil,
		Path:     "/AAA",
		Headers: map[string][]string{
			"TestHeaders":  {"123"},
			"Content-Type": {"application/json"},
		},
		Body: result,
		Root: false,
	}
	resp := api.Fetch(req)
	fmt.Printf("%+v", resp)

	outdoor := plugin.resource.GetOutDoor()

	req2 := &common_type.HttpRequest{
		Method:   "get",
		QueryMap: nil,
		Url:      "http://localhost:9001/anhao",
		Headers: map[string][]string{
			"TestHeaders":  {"123"},
			"Content-Type": {"application/json"},
		},
		Body: result,
		Root: false,
	}
	outdoor.AsyncFetch(req2, func(response *common_type.HttpResponse, pluginError common_type.PluginError) {
		fmt.Printf("%s%s", string(response.Body), response.Err)
	})

	time.Sleep(time.Hour)
}

func testDB(plugin *mockPlugin) {
	localDB := plugin.resource.GetLocalDB()

	//if err := localDB.ImportSQL("./config/plugin.sql"); err != nil {
	//	panic(err)
	//}
	//
	//rawData, cols, err := localDB.Select("select * from upload_tips;")
	//if err != nil {
	//	panic(err)
	//}
	//type User struct {
	//	UUID       string `orm:"team_uuid"`
	//	Name       string `orm:"content"`
	//	UpdateTime int    `orm:"update_time"`
	//}
	//users := make([]*User, 0)
	//e := localDB.Unmarshal(rawData, cols, &users)
	//if e != nil {
	//	panic(e)
	//}
	//fmt.Printf("%+v\n", users)

	localDB.AsyncSelect("select * from upload_tips;", func(data []*common_type.RawData, descs []*common_type.ColumnDesc, pluginError common_type.PluginError) {
		type User struct {
			UUID       string `orm:"team_uuid"`
			Name       string `orm:"content"`
			UpdateTime int    `orm:"update_time"`
		}
		users := make([]*User, 0)
		e := localDB.Unmarshal(data, descs, &users)
		if e != nil {
			panic(e)
		}
		fmt.Printf("%+v\n", users)
	})

	//err := localDB.Exec("update upload_tips set content = 'yyy' where team_uuid = 'asd';")
	//if err != nil {
	//	panic(err)
	//}

	sysDB := plugin.resource.GetSysDB()
	rawData, cols, err := sysDB.Select("plugin_platform", "select * from plugin_package;")
	if err != nil {
		panic(err)
	}

	type Package struct {
		ID         string `orm:"id"`
		AppUUID    string `orm:"app_uuid"`
		Name       string `orm:"name"`
		Size       int    `orm:"size"`
		Version    string `orm:"version"`
		CreateTime int    `orm:"create_time"`
		UpdateTime int    `orm:"update_time"`
		Deleted    int    `orm:"deleted"`
	}
	ps := make([]*Package, 0)
	e := sysDB.Unmarshal(rawData, cols, &ps)
	if e != nil {
		panic(e)
	}
	fmt.Printf("%+v\n", ps)
	time.Sleep(time.Hour)
}

func testWorkSpace(plugin *mockPlugin) {
	workSpace := plugin.resource.GetWorkspace()
	workSpace.AsyncCopy("qwe_____.txt", "qwe_____2.txt", func(pluginError common_type.PluginError) {
		fmt.Println("----___ ok ___----")
		fmt.Println(pluginError)
	})
	time.Sleep(time.Hour)
	e := workSpace.CreateFile("qwe_____.txt")
	if e != nil {
		panic(e)
	}
	e = workSpace.WriteStrings("qwe_____.txt", []string{"123\n", "345"})
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

func (i *mockPluginDescriptor) ApplicationID() string { return "lt1ZZuMd" }
func (i *mockPluginDescriptor) Name() string          { return "Application123" }
func (i *mockPluginDescriptor) Language() string      { return "Language123" }
func (i *mockPluginDescriptor) LanguageVersion() common_type.IVersion {
	return common_type.NewVersion(1, 2, 3)
}
func (i *mockPluginDescriptor) ApplicationVersion() common_type.IVersion {
	return common_type.NewVersion(1, 0, 0)
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
