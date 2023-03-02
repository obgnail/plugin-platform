package main

import (
	"github.com/obgnail/plugin-platform/host_boot/handler"
	"time"
)

func main() {
	h := handler.Default()
	h.Run()

	time.Sleep(time.Hour)
}
