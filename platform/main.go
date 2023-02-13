package main

import (
	"fmt"
	"github.com/obgnail/plugin-platform/platform/config"
	"github.com/obgnail/plugin-platform/platform/model/mysql"
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
	onStart(mysql.InitDB)
}

func main() {
	Init()
	router.Run()
}
