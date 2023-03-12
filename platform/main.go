package main

import (
	"fmt"
	hotboot_handler "github.com/obgnail/plugin-platform/host_boot/handler"
	"github.com/obgnail/plugin-platform/platform/conn/handler"
	"github.com/obgnail/plugin-platform/platform/conn/hub"
	"github.com/obgnail/plugin-platform/platform/model/mysql"
	"github.com/obgnail/plugin-platform/platform/model/redis"
	"github.com/obgnail/plugin-platform/platform/router"
)

func main() {
	Init()
	router.Run()
}

func Init() {
	OnStart(handler.InitPlatformHandler)
	OnStart(mysql.InitDB)
	OnStart(redis.InitRedis)
	OnStart(hotboot_handler.InitHostBoot)
	OnStart(hub.InitHub)
}

func OnStart(fn func() error) {
	if err := fn(); err != nil {
		panic(fmt.Sprintf("Error at onStart: %s\n", err))
	}
}
