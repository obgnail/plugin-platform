package main

import (
	"fmt"
	"github.com/obgnail/plugin-platform/common/log"
	"github.com/obgnail/plugin-platform/platform/handler/handler"
	"github.com/obgnail/plugin-platform/platform/model/mysql"
	"github.com/obgnail/plugin-platform/platform/model/redis"
	"github.com/obgnail/plugin-platform/platform/pool/plugin_pool"
	"time"
)

func main() {
	Init()
	h := handler.Default()
	h.Run()

	//h.Send(&protocol.PlatformMessage{}, 30*time.Second)

	log.Info("PlatformHandler OK")

	go func() {
		time.Sleep(15 * time.Second)
		log.Info("InstallPlugin...")
		h.InstallPlugin("lt1ZZuMd", "InstanceID123", "上传文件的安全提示",
			"golang", "1.14.0", "1.0.0")
		h.StartPlugin("lt1ZZuMd", "InstanceID123", "上传文件的安全提示",
			"golang", "1.14.0", "1.0.0")

		time.Sleep(time.Second * 20)
		log.Info("kill Plugin")
		h.KillPlugin("InstanceID123")
	}()
	time.Sleep(time.Hour)
}

func onStart(fn func() error) {
	if err := fn(); err != nil {
		panic(fmt.Sprintf("Error at onStart: %s\n", err))
	}
}

func Init() {
	onStart(plugin_pool.InitPluginPool)
	onStart(mysql.InitDB)
	onStart(redis.InitRedis)
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
