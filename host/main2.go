package main

import (
	"encoding/json"
	"fmt"
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/log"
	"github.com/obgnail/plugin-platform/host/resource/release"
	"reflect"
	"runtime"
	"time"
)

func main() {
	plugin := &mockPlugin2{}
	resource := &release.Resource{Plugin: plugin, Sender: nil}
	instanceDesc := &mockInstanceDesc2{}

	log.Info("OK")
	err := plugin.Assign(instanceDesc, resource)
	if err != nil {
		panic(err)
	}

	var w common_type.IPlugin = &Wrapper{IPlugin: plugin}

	h := &common_type.HttpRequest{Url: "asdasdasdasdasdULR"}
	result := invoke(w, "OnHttpCall", h)
	aaa := result[0].Interface().(*common_type.HttpResponse)
	fmt.Println(aaa)

	//testNetwork(plugin)
}

type Wrapper struct {
	common_type.IPlugin
}

func invoke(any interface{}, name string, args ...interface{}) []reflect.Value {
	defer func() {
		if err := recover(); err != nil {
			buf := make([]byte, 10240)
			n := runtime.Stack(buf, false)
			stackInfo := fmt.Sprintf("%s", buf[:n])
			er := fmt.Errorf("plugin invoke panic: %v, %s", err, stackInfo)
			panic(er)
		}
	}()

	inputs := make([]reflect.Value, len(args))
	for i := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	ptr := reflect.ValueOf(any).Elem().MethodByName(name)
	data := ptr.Call(inputs)
	return data
}

func testNetwork(plugin *mockPlugin2) {
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

////////////////////////////
type mockInstanceDesc2 struct{}

func (i *mockInstanceDesc2) InstanceID() string { return "InstanceID123" }
func (i *mockInstanceDesc2) PluginDescription() common_type.IPluginDescriptor {
	return &mockPluginDescriptor2{}
}

type mockPluginDescriptor2 struct{}

func (i *mockPluginDescriptor2) ApplicationID() string { return "lt1ZZuMd" }
func (i *mockPluginDescriptor2) Name() string          { return "Application123" }
func (i *mockPluginDescriptor2) Language() string      { return "Language123" }
func (i *mockPluginDescriptor2) LanguageVersion() common_type.IVersion {
	return common_type.NewVersion(1, 2, 3)
}
func (i *mockPluginDescriptor2) ApplicationVersion() common_type.IVersion {
	return common_type.NewVersion(1, 0, 0)
}
func (i *mockPluginDescriptor2) HostVersion() common_type.IVersion {
	return common_type.NewVersion(1, 2, 3)
}
func (i *mockPluginDescriptor2) MinSystemVersion() common_type.IVersion {
	return common_type.NewVersion(1, 2, 3)
}

////////////////////////////////////////////////////////////////////////////////////////////////
var _ common_type.IPlugin = (*mockPlugin2)(nil)

type mockPlugin2 struct {
	descriptor common_type.IInstanceDescription
	resource   common_type.IResources
}

func GetPlugin() common_type.IPlugin {
	return &mockPlugin2{}
}
func (p *mockPlugin2) GetPluginDescription() common_type.IInstanceDescription {
	return p.descriptor
}

func (p *mockPlugin2) Assign(pd common_type.IInstanceDescription, resource common_type.IResources) common_type.PluginError {
	p.descriptor = pd
	p.resource = resource
	return nil
}

func (p *mockPlugin2) OnHttpCall(req *common_type.HttpRequest) (resp *common_type.HttpResponse) {
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

func (p *mockPlugin2) Enable(common_type.LifeCycleRequest) common_type.PluginError {
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

func (p *mockPlugin2) Disable(common_type.LifeCycleRequest) common_type.PluginError {
	err := p.resource.GetWorkspace().CreateFile("mack_plugin_test.txt")
	if err != nil {
		log.ErrorDetails(err)
	}
	return nil
}

func (p *mockPlugin2) Start(common_type.LifeCycleRequest) common_type.PluginError {
	err := p.resource.GetWorkspace().WriteBytes("mack_plugin_test.txt", []byte("1234567654321"))
	if err != nil {
		log.ErrorDetails(err)
	}
	return nil
}

func (p *mockPlugin2) Stop(common_type.LifeCycleRequest) common_type.PluginError {
	return nil
}

// TODO
func (p *mockPlugin2) CheckState() common_type.PluginError {
	return nil
}

// TODO
func (p *mockPlugin2) CheckCompatibility() common_type.PluginError {
	return nil
}

func (p *mockPlugin2) Install(common_type.LifeCycleRequest) common_type.PluginError {
	return nil
}

func (p *mockPlugin2) UnInstall(common_type.LifeCycleRequest) common_type.PluginError {
	return nil
}

func (p *mockPlugin2) Upgrade(common_type.IVersion, common_type.LifeCycleRequest) common_type.PluginError {
	return nil
}

func (p *mockPlugin2) OnEvent(eventType string, payload interface{}) common_type.PluginError {
	return nil
}

func (p *mockPlugin2) OnExternalHttpRequest(request *common_type.HttpRequest) *common_type.HttpResponse {
	return nil
}
