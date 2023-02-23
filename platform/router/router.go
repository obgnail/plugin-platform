package router

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/obgnail/plugin-platform/common/log"
	"github.com/obgnail/plugin-platform/platform/config"
	"github.com/obgnail/plugin-platform/platform/controllers"
	utils_errors "github.com/obgnail/plugin-platform/utils/errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run() {
	gin.SetMode(gin.ReleaseMode)
	app := gin.Default()

	plugin := app.Group("/plugin")

	plugin.POST("/list", controllers.ListPlugins)
	// 插件的路由同步给标品
	plugin.GET("router_list", controllers.RouterList)

	// life cycle
	plugin.POST("/upload", controllers.Upload)
	plugin.POST("/install", controllers.Install)
	//plugin.POST("/delete_opk", controllers.DeleteOpk)
	//plugin.POST("/uninstall", controllers.UnInstall)
	//plugin.POST("/enable", controllers.Enable)
	//plugin.POST("/disable", controllers.Disable)
	//plugin.POST("/upgrade", controllers.Upgrade)

	addr := fmt.Sprintf("%s:%d", config.StringOrPanic("host"), config.IntOrPanic("http_port"))
	srv := &http.Server{Addr: addr, Handler: app}
	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.ErrorDetails(utils_errors.Trace(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown: %v", err)
	}
	//common.RemoveAllHosts() TODO

	log.Info("Server exiting")
}
