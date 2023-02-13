package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/obgnail/plugin-platform/platform/service/lifecycle"
)

func Upload(c *gin.Context) {
	resp, result := lifecycle.Upload(c)
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

//func DeleteOpk(c *gin.Context) {
//	req := &lifecycle_action.DeleteOpkReq{}
//	if err := c.BindJSON(&req); err != nil {
//		return
//	}
//	resp, result := lifecycle_action.DeleteOpk(req)
//	RenderJSON(c, result, resp)
//}
//

//
//func UnInstall(c *gin.Context) {
//	req := &lifecycle_action.UninstallReq{}
//	if err := c.BindJSON(&req); err != nil {
//		return
//	}
//	userUUID := c.GetHeader("Ones-User-Id")
//	userToken := c.GetHeader("Ones-Auth-Token")
//	resp, result := lifecycle_action.Uninstall(req, userUUID, userToken)
//	RenderJSON(c, result, resp)
//}
//
//func Enable(c *gin.Context) {
//	req := &lifecycle_action.EnableReq{}
//	if err := c.BindJSON(&req); err != nil {
//		return
//	}
//	userID := c.GetHeader("Ones-User-Id")
//	resp, result := lifecycle_action.Enable(req, userID)
//	RenderJSON(c, result, resp)
//}
//
//func Disable(c *gin.Context) {
//	req := &lifecycle_action.DisableReq{}
//	if err := c.BindJSON(&req); err != nil {
//		return
//	}
//	userID := c.GetHeader("Ones-User-Id")
//	resp, result := lifecycle_action.Disable(req, userID)
//	RenderJSON(c, result, resp)
//}
//
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
