package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/obgnail/plugin-platform/platform/service/lifecycle"
)

func Upload(c *gin.Context) {
	resp, result := lifecycle.Upload(c)
	RenderJSON(c, result, resp)
}

func Delete(c *gin.Context) {
	req := &lifecycle.DeleteReq{}
	if err := c.BindJSON(&req); err != nil {
		return
	}
	resp, result := lifecycle.Delete(req)
	RenderJSON(c, result, resp)
}

func Install(c *gin.Context) {
	req := &lifecycle.InstallReq{}
	if err := c.BindJSON(&req); err != nil {
		return
	}
	resp, result := lifecycle.Install(req)
	RenderJSON(c, result, resp)
}

func UnInstall(c *gin.Context) {
	req := &lifecycle.UninstallReq{}
	if err := c.BindJSON(&req); err != nil {
		return
	}
	resp, result := lifecycle.Uninstall(req)
	RenderJSON(c, result, resp)
}

func Enable(c *gin.Context) {
	req := &lifecycle.EnableReq{}
	if err := c.BindJSON(&req); err != nil {
		return
	}
	resp, result := lifecycle.Enable(req)
	RenderJSON(c, result, resp)
}

func Disable(c *gin.Context) {
	req := &lifecycle.DisableReq{}
	if err := c.BindJSON(&req); err != nil {
		return
	}
	resp, result := lifecycle.Disable(req)
	RenderJSON(c, result, resp)
}

//func Upgrade(c *gin.Context) {
//	req := &lifecycle_action.UpgradeReq{}
//	if err := c.BindJSON(&req); err != nil {
//		return
//	}
//	resp, result := lifecycle_action.Upgrade(req)
//	RenderJSON(c, result, resp)
//}
//
//type UpgradeReq struct {
//	AppUUID  string `json:"app_uuid"`
//	OrgUUID  string `json:"organization_uuid"`
//	TeamUUID string `json:"team_uuid"`
//}
