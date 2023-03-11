package hub

import (
	"fmt"
	"github.com/obgnail/plugin-platform/common/config"
	"github.com/obgnail/plugin-platform/common/log"
	"time"
)

var (
	cycleInterval = time.Duration(config.Int("main_system.heartbeat_sec", 3)) * time.Second
	addr          = fmt.Sprintf("http://%s:%d",
		config.String("platform.host", "127.0.0.1"),
		config.Int("platform.http_port", 9005),
	)

	routerListPath = "/plugin/router_list"
)

func InitPluginService() {
	go func() {
		ticker := time.NewTicker(cycleInterval)
		defer ticker.Stop()

		url := addr + routerListPath
		for {
			select {
			case <-ticker.C:
				resp, err := Get(url, nil)
				if err != nil {
					log.ErrorDetails(err)
					continue
				}
				plugins, err := unmarshalPlugins(resp)
				if err != nil {
					log.ErrorDetails(err)
					continue
				}
				if err = registerHub(plugins); err != nil {
					log.ErrorDetails(err)
					continue
				}
			}
		}
	}()
}
