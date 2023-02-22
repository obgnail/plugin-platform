package main

import (
	"fmt"
	"github.com/obgnail/plugin-platform/platform/config"
	"github.com/obgnail/plugin-platform/platform/model/mysql"
	"github.com/obgnail/plugin-platform/platform/pool/plugin_pool"
	"github.com/obgnail/plugin-platform/platform/router"
	"github.com/obgnail/plugin-platform/utils/log"
)

func onStart(fn func() error) {
	if err := fn(); err != nil {
		panic(fmt.Sprintf("Error at onStart: %s\n", err))
	}
}

func Init() {
	onStart(config.InitConfig)
	onStart(log.InitLogger)
	onStart(plugin_pool.InitPluginPool)
	onStart(mysql.InitDB)
}

func RunServer() {
	go plugin_pool.Run()
	router.Run()
}

func main() {
	Init()
	RunServer()
}
