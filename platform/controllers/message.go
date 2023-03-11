package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/obgnail/plugin-platform/platform/service/message"
)

func ListPlugins(c *gin.Context) {
	resp, result := message.ListPlugins()
	RenderJSON(c, result, resp)
}

func RouterList(c *gin.Context) {
	resp, result := message.RouterList()
	RenderJSON(c, result, resp)
}
