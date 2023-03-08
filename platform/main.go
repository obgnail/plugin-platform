package main

import (
	"fmt"
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/log"
	hotboot_handler "github.com/obgnail/plugin-platform/host_boot/handler"
	"github.com/obgnail/plugin-platform/platform/conn/handler"
	"github.com/obgnail/plugin-platform/platform/model/mysql"
	"github.com/obgnail/plugin-platform/platform/model/redis"
	"github.com/obgnail/plugin-platform/platform/router"
	"time"
)

func main() {
	Init()
	log.Info("run")

	go func() {
		time.Sleep(4 * time.Second)
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
		result := <-handler.CallPluginInternalHTTP("InstanceID123", &req, "OnHttpCall")
		fmt.Printf("==============Internal=================%+v\n", result)

		result = <-handler.CallPluginExternalHTTP("InstanceID123", &req)
		fmt.Printf("============External===================%+v\n", result)

		err := <-handler.CallPluginEvent("InstanceID123", "myEventType", []byte("xasasdasdasdasd"))
		fmt.Printf("((((((((((((((((((((((((((((((((((%+v\n", err.Msg())

		err = <-handler.CallPluginConfigChanged("InstanceID123", "myConfigKey", []string{"originValue"}, []string{"newValue"})
		fmt.Printf("``````````````````````````````````%+v\n", err.Msg())

		resp := <-handler.CallPluginFunction("InstanceID123", "abilityID", "abilityType", "AbilityFuncKey1", []byte("args1"))
		fmt.Println("+++++++++++++++1111111111+++++++++++++++++++++++", resp.Data, resp.Err)
		resp = <-handler.CallPluginFunction("InstanceID123", "abilityID", "abilityType", "AbilityFuncKey2", []byte("args2"))
		fmt.Println("+++++++++++++++2222222222+++++++++++++++++++++++", resp.Data, resp.Err)
	}()

	router.Run()
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
