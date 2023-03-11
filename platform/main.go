package main

import (
	"fmt"
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/errors"
	"github.com/obgnail/plugin-platform/common/log"
	hotboot_handler "github.com/obgnail/plugin-platform/host_boot/handler"
	"github.com/obgnail/plugin-platform/platform/conn/handler"
	"github.com/obgnail/plugin-platform/platform/conn/hub"
	"github.com/obgnail/plugin-platform/platform/conn/hub/ability"
	hub_router "github.com/obgnail/plugin-platform/platform/conn/hub/router"
	"github.com/obgnail/plugin-platform/platform/model/mysql"
	"github.com/obgnail/plugin-platform/platform/model/redis"
	"github.com/obgnail/plugin-platform/platform/router"
	"github.com/obgnail/plugin-platform/platform/service/common"
	"github.com/obgnail/plugin-platform/platform/service/utils"
	"gopkg.in/yaml.v2"
	"time"
)

func main() {
	Init()
	log.Info("run")

	cfg, err := LoadYamlConfig("lt1ZZuMd", "1.0.0")
	if err != nil {
		panic(err)
	}

	if err := hub_router.RegisterRouter("HXCEB1oF", cfg.Apis); err != nil {
		panic(err)
	}
	ability.RegisterAbility("HXCEB1oF", cfg.Abilities)

	go func() {
		time.Sleep(4 * time.Second)
		log.Info("InstallPlugin...")
		<-handler.InstallPlugin("lt1ZZuMd", "HXCEB1oF", "上传文件的安全提示",
			"golang", "1.14.0", "1.0.0")
		<-handler.EnablePlugin("HXCEB1oF")

		req := common_type.HttpRequest{
			Method:   "",
			QueryMap: nil,
			Url:      "/urlXXXXXXXXXXXXX",
			Path:     "",
			Headers:  nil,
			Body:     nil,
			Root:     false,
		}
		result := <-handler.CallPluginExternalHTTP("HXCEB1oF", &req)
		fmt.Printf("==============Internal=================%+v\n", result)

		//err := <-handler.CallPluginEvent("InstanceID123", "project.task", []byte("project.task_payload"))
		//fmt.Printf("111%+v\n", err.Msg())
		//err = <-handler.CallPluginEvent("InstanceID123", "project.user", []byte("project.user_payload"))
		//fmt.Printf("222%+v\n", err.Msg())
		//err = <-handler.CallPluginEvent("InstanceID123", "project.userXXX", []byte("project.user_payload"))
		//fmt.Printf("333%+v\n", err.Msg())
		//
		//req := common_type.HttpRequest{
		//	Method:   "",
		//	QueryMap: nil,
		//	Url:      "/url",
		//	Path:     "",
		//	Headers:  nil,
		//	Body:     nil,
		//	Root:     false,
		//}
		//result := <-handler.CallPluginInternalHTTP("InstanceID123", &req, "OnHttpCall")
		//fmt.Printf("==============Internal=================%+v\n", result)
		//
		//result = <-handler.CallPluginExternalHTTP("InstanceID123", &req)
		//fmt.Printf("============External===================%+v\n", result)
		//
		//err := <-handler.CallPluginEvent("InstanceID123", "myEventType", []byte("xasasdasdasdasd"))
		//fmt.Printf("((((((((((((((((((((((((((((((((((%+v\n", err.Msg())
		//
		//err = <-handler.CallPluginConfigChanged("InstanceID123", "myConfigKey", []string{"originValue"}, []string{"newValue"})
		//fmt.Printf("``````````````````````````````````%+v\n", err.Msg())
		//
		//resp := <-handler.CallPluginFunction("InstanceID123", "abilityID", "abilityType", "AbilityFuncKey1", []byte("args1"))
		//fmt.Println("+++++++++++++++1111111111+++++++++++++++++++++++", resp.Data, resp.Err)
		//resp = <-handler.CallPluginFunction("InstanceID123", "abilityID", "abilityType", "AbilityFuncKey2", []byte("args2"))
		//fmt.Println("+++++++++++++++2222222222+++++++++++++++++++++++", resp.Data, resp.Err)
	}()

	router.Run()
}

func LoadYamlConfig(appid, version string) (*common.PluginConfig, error) {
	yamlPath := utils.GetPluginConfigPath(appid, version)
	res, err := utils.ReadFile(yamlPath)
	if err != nil {
		return nil, errors.Trace(err)
	}

	var pluginConfig = new(common.PluginConfig)
	if err := yaml.Unmarshal(res, pluginConfig); err != nil {
		return nil, errors.Trace(err)
	}
	return pluginConfig, nil
}
func main2() {
	Init()

	//h.Send(&protocol.PlatformMessage{}, 30*time.Second)

	log.Info("PlatformHandler OK")

	go func() {
		time.Sleep(15 * time.Second)
		log.Info("InstallPlugin...")
		<-handler.InstallPlugin("lt1ZZuMd", "InstanceID123", "上传文件的安全提示",
			"golang", "1.14.0", "1.0.0")
		<-handler.EnablePlugin("InstanceID123")

		req := common_type.HttpRequest{
			Method:   "",
			QueryMap: nil,
			Url:      "/url",
			Path:     "",
			Headers:  nil,
			Body:     nil,
			Root:     false,
		}
		fmt.Printf("+++++++++++++++++++++++++%+v\n", req)
		result := <-handler.CallPluginHTTP("InstanceID123", &req, true, "OnHttpCall")
		fmt.Printf("===============================%+v\n", result)
		//time.Sleep(time.Second * 20)
		//log.Info("kill Plugin")
		//handler.KillPlugin("InstanceID123")
	}()
	time.Sleep(time.Hour)
}

func onStart(fn func() error) {
	if err := fn(); err != nil {
		panic(fmt.Sprintf("Error at onStart: %s\n", err))
	}
}

func Init() {
	onStart(handler.InitPlatformHandler)
	onStart(mysql.InitDB)
	onStart(redis.InitRedis)
	onStart(hotboot_handler.InitHostBoot)
	onStart(hub.InitHub)
}

//func RunServer() {
//	go plugin_pool.Run()
//	router.Run()
//}
//
//func main2() {
//	Init()
//	RunServer()
//}
