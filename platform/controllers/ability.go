package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/obgnail/plugin-platform/platform/conn/hub/ability"
	"strconv"
)

func CallAbility(c *gin.Context) {
	args, err := c.GetRawData()
	if err != nil {
		RenderError(c, err)
		return
	}
	instanceID := c.Request.Header.Get("instanceID")
	abilityID := c.Request.Header.Get("abilityID")
	abilityType := c.Request.Header.Get("abilityType")
	abilityFunc := c.Request.Header.Get("abilityFunc")

	data, err := ability.SyncExecute(instanceID, abilityID, abilityType, abilityFunc, args)
	if err != nil {
		RenderError(c, err)
		return
	}
	c.Header("Content-Length", strconv.Itoa(len(data)))
	c.Render(200, render.Data{Data: data})
}
