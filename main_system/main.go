package main

import (
	"github.com/obgnail/plugin-platform/main_system/platform/hub"
	"github.com/obgnail/plugin-platform/main_system/router"
)

func main() {
	Init()
	router.Run()
}

func Init() {
	hub.InitPluginService()
}
