package router

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/obgnail/plugin-platform/common/config"
	utils_errors "github.com/obgnail/plugin-platform/common/errors"
	"github.com/obgnail/plugin-platform/common/log"
	"github.com/obgnail/plugin-platform/main_system/middlewares"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run() {
	gin.SetMode(gin.ReleaseMode)
	app := gin.Default()

	registerPluginMiddlewares(app)
	registerRouter(app)
	run(app)
}

func registerPluginMiddlewares(app *gin.Engine) {
	if !config.Bool("main_system.enable_plugin", true) {
		return
	}

	app.Use(middlewares.PrefixProcessor()) // 顺序不能反,先prefix再replace
	app.Use(middlewares.ReplaceProcessor())
	app.Use(middlewares.SuffixProcessor())
	app.NoRoute(middlewares.AdditionProcessor())
}

func run(app *gin.Engine) {
	addr := config.StringOrPanic("main_system.addr")
	srv := &http.Server{Addr: addr, Handler: app}
	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.ErrorDetails(utils_errors.Trace(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Warn("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown: %v", err)
	}

	log.Warn("Server exiting")
}
