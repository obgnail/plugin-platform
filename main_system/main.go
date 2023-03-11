package main

import (
	"github.com/obgnail/plugin-platform/common/log"
	"github.com/obgnail/plugin-platform/main_system/platform/hub"
	"github.com/obgnail/plugin-platform/main_system/router"
)

func main() {
	log.Info("run main system")
	Init()

	router.Run()
}

func Init() {
	hub.InitPluginService()
}
