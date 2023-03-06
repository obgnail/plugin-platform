package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/obgnail/plugin-platform/common/errors"
	"github.com/obgnail/plugin-platform/platform/controllers"
	"strings"
)

const (
	PluginPrefix = "plugin_external"
)

func PluginInvoke() gin.HandlerFunc {
	return invoke
}

func invoke(c *gin.Context) {
	path := c.Request.RequestURI
	parts := strings.Split(path, PluginPrefix)
	if len(parts) < 2 {
		handlerError(c, "path error")
		return
	}
	//url := parts[1]
	//body, err := c.GetRawData()
	//if err != nil {
	//	handlerError(c, err.Error())
	//	return
	//}

	//instanceUUID = plugin_api.GetInstanceUUIDByPathMethod(url, c.Request.Method)
}

func handlerError(c *gin.Context, reason string) {
	err := errors.PluginCallError(errors.CallPluginFailed, reason)
	controllers.RenderJSONAndStop(c, err, nil)
}
