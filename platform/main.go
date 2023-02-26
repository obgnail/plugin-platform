package main

import (
	"fmt"
	"github.com/obgnail/plugin-platform/common/log"
	"github.com/obgnail/plugin-platform/platform/handler"
	"github.com/obgnail/plugin-platform/platform/model/mysql"
	"github.com/obgnail/plugin-platform/platform/model/redis"
	"github.com/obgnail/plugin-platform/platform/pool/plugin_pool"
	"time"
)

func main() {
	Init()
	h := handler.Default()
	h.Run()

	log.Info("PlatformHandler OK")

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
