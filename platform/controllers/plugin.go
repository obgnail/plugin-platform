package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/obgnail/plugin-platform/platform/service/plugin"
)

func ListPlugins(c *gin.Context) {
	resp, result := plugin.ListPlugins()
	RenderJSON(c, result, resp)
}

func RouterList(c *gin.Context) {
	//resp := plugin_data.RouterList()
	//RenderJSON(c, nil, resp)
}
