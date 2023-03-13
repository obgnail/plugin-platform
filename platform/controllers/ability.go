package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/obgnail/plugin-platform/platform/conn/handler"
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
	c.Header("Content-Type", "application/json")
	c.Header("Content-Length", strconv.Itoa(len(data)))
	c.Render(200, render.Data{Data: data})
}

func OnEvent(c *gin.Context) {
	payload, err := c.GetRawData()
	if err != nil {
		RenderError(c, err)
		return
	}

	instanceID := c.Request.Header.Get("instanceID")
	eventType := c.Request.Header.Get("eventType")
	pluginErr := <-handler.CallPluginEvent(instanceID, eventType, payload)
	if pluginErr != nil {
		e := fmt.Errorf(pluginErr.Msg())
		RenderError(c, e)
		return
	}
	RenderJSON(c, nil, nil)
	return
}
