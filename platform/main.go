package main

import (
	"fmt"
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

	//go func() {
	//	time.Sleep(15 * time.Second)
	//	log.Info("InstallPlugin...")
	//	handler.InstallPlugin("lt1ZZuMd", "InstanceID123", "上传文件的安全提示",
	//		"golang", "1.14.0", "1.0.0")
	//	handler.EnablePlugin("lt1ZZuMd", "InstanceID123", "上传文件的安全提示",
	//		"golang", "1.14.0", "1.0.0")
	//
	//	//time.Sleep(time.Second * 20)
	//	//log.Info("kill Plugin")
	//	//handler.KillPlugin("InstanceID123")
	//}()

	router.Run()
}

func main2() {
	Init()

	//h.Send(&protocol.PlatformMessage{}, 30*time.Second)

	log.Info("PlatformHandler OK")

	go func() {
		time.Sleep(15 * time.Second)
		log.Info("InstallPlugin...")
		handler.InstallPlugin("lt1ZZuMd", "InstanceID123", "上传文件的安全提示",
			"golang", "1.14.0", "1.0.0")
		handler.EnablePlugin("lt1ZZuMd", "InstanceID123", "上传文件的安全提示",
			"golang", "1.14.0", "1.0.0")

		time.Sleep(time.Second * 20)
		log.Info("kill Plugin")
		handler.KillPlugin("InstanceID123")
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
