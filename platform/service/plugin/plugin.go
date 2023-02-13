package plugin

//
//import (
//	"encoding/json"
//	"github.com/BangWork/ones-platform-api/app/pool/plugin"
//)
//
//type Output struct {
//	PluginInfoList  []*PluginInfo `json:"plugin_info_list"`
//	TimeStamp       int64         `json:"server_update_stamp"`
//	AbilityInfoList AbilityInfo   `json:"ability_info_list"`
//	AbilityApis     []*Router     `json:"ability_apis"`
//}
//
//func RouterList() (output Output) {
//	instanceModel := mysql.ModelPluginInstance()
//
//	var instances = make([]*mysql.PluginInstance, 0)
//	arg := &mysql.PluginInstance{}
//	err := instanceModel.All(&instances, arg)
//	if err != nil {
//		log.Logger.Error("instanceModel.All err: %v", err)
//		return Output{}
//	}
//
//	pluginInfoList := make([]*PluginInfo, 0)
//	for _, l := range instances {
//		var apis = make([]*plugin.Api, 0)
//		_ = json.Unmarshal([]byte(l.Apis), &apis)
//
//		uModel := mysql.ModelPluginUser()
//		u := &mysql.PluginUser{
//			OrgUUID:      l.OrgUUID,
//			TeamUUID:     l.TeamUUID,
//			AppUUID:      l.AppUUID,
//			InstanceUUID: l.InstanceUUID,
//		}
//		_ = uModel.One(u)
//
//		userTeamList := make([]*UserTeam, 0)
//		userTeamList = append(userTeamList, &UserTeam{
//			TeamUUID: l.TeamUUID,
//			UserUUID: u.UserUUID,
//		})
//
//		routerList := make([]*Router, 0)
//		for _, api := range apis {
//			r := &Router{
//				Name:   "",
//				Type:   api.Type,
//				Method: api.Methods,
//				Url:    api.Url,
//				Scope:  api.Scope,
//			}
//			routerList = append(routerList, r)
//		}
//
//		pluginStatus, _ := plugin.PluginPool.GetPluginStatus(l.InstanceUUID)
//		if pluginStatus == -1 {
//			pluginStatus = l.Status
//		}
//		p := &PluginInfo{
//			Instance: &Instance{
//				UUID:         l.InstanceUUID,
//				Name:         l.Name,
//				Type:         1, // 团队插件
//				LifeStage:    pluginStatus,
//				OrgUUID:      l.OrgUUID,
//				TeamUUID:     l.TeamUUID,
//				UserTeamList: userTeamList,
//				Version:      l.Version,
//				Description:  l.Description,
//			},
//			Routers: routerList,
//		}
//		pluginInfoList = append(pluginInfoList, p)
//	}
//
//	output.AbilityApis = getPlatformAbilityRouter() // 平台能力路由
//
//	// 获取能力表对应的业务表的关联ID
//	list, _ := Resolver.GetAbilityRelationIDList()
//	output.AbilityInfoList = AbilityInfo{
//		AbilityRelationInfo: list,
//	}
//	output.PluginInfoList = pluginInfoList
//	// TODO
//	output.TimeStamp = plugin.PluginPool.TimeStamp
//
//	return output
//}
//
//// appendPlatformAbilityRouter 给能力注入一个 route 到 bang-api
//func getPlatformAbilityRouter() []*Router {
//	apis := make([]*Router, 0)
//	if true {
//		// todo 判断是否有插件使用了这个能力
//		apis = append(apis, &Router{
//			Name:   "SimpleAuthLogin",
//			Type:   "ability",
//			Method: []string{"POST"},
//			Url:    "/ability/simple_auth_login",
//		})
//	}
//	return apis
//}
