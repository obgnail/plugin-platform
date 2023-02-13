package plugin

//
//import (
//	"encoding/json"
//	"github.com/BangWork/ones-platform-api/app/pool/plugin"
//	"github.com/BangWork/ones-platform-api/app/service/common/i18n"
//	"github.com/gin-gonic/gin"
//	"github.com/obgnail/plugin-platform/platform/model/mysql"
//	"gopkg.in/yaml.v2"
//	"os"
//	"sort"
//)
//
//func ListPlugins() (ret gin.H, err error) {
//	var instances = make([]*mysql.PluginInstance, 0)
//
//	instanceModel := mysql.ModelPluginInstance()
//	var instanceArg = &mysql.PluginInstance{
//		OrgUUID:  req.OrgUUID,
//		TeamUUID: req.TeamUUID,
//	}
//	err = instanceModel.All(&instances, instanceArg)
//	if err != nil {
//		err = errors.Wrapf(err, errors.ServerError, "model.ModelPluginInstance.All err: %v", err)
//		return
//	}
//
//	packageModel := mysql.ModelPluginPackage()
//	var packages = make([]*mysql.PluginPackage, 0)
//	var packageArg = &mysql.PluginPackage{
//		OrgUUID:  req.OrgUUID,
//		TeamUUID: req.TeamUUID,
//	}
//	if err = packageModel.All(&packages, packageArg); err != nil {
//		err = errors.Wrapf(err, errors.ServerError, "model.ModelPluginPackage.All err: %v", err)
//		return
//	}
//	var packageInfo = make(map[string]*mysql.PluginPackage)
//	for _, v := range packages {
//		index := v.AppUUID + v.Version
//		packageInfo[index] = v
//	}
//
//	var webFilePath string
//	var iconPath string
//
//	all := make([]*plugin.PluginConfig, 0)
//	for _, instance := range instances {
//		var apis = make([]*plugin.Api, 0)
//		var modules = make([]interface{}, 0)
//		var permissions = make([]*plugin.Permission, 0)
//		var abilities = make([]*plugin.Ability, 0)
//
//		_ = json.Unmarshal([]byte(instance.Apis), &apis)
//		_ = json.Unmarshal([]byte(instance.Modules), &modules)
//		_ = json.Unmarshal([]byte(instance.Abilities), &abilities)
//		index := instance.AppUUID + instance.Version
//
//		pluginPackageInfo, ok := packageInfo[index]
//		if ok {
//			webFilePath = pluginPackageInfo.WebFilePath
//			if pluginPackageInfo.Icon != "" {
//				iconPath = string(os.PathSeparator) + "plugin" + webFilePath + "logo.svg"
//			} else {
//				iconPath = pluginPackageInfo.Icon
//			}
//
//			configYaml := appservice.LoadYamlFile(pluginPackageInfo.ConfigFilePath)
//			if configYaml == "" {
//				err = errors.Wrapf(err, errors.ServerError, "loadYamlFile failed")
//				return
//			}
//
//			var pluginConfig = new(plugin.PluginConfig)
//			if err = yaml.Unmarshal([]byte(configYaml), pluginConfig); err != nil {
//				return
//			}
//			permissions = pluginConfig.Permission
//		}
//
//		var service = &plugin.Service{
//			AppUUID:      instance.AppUUID,
//			InstanceUUID: instance.InstanceUUID,
//			Name:         instance.Name,
//			Version:      instance.Version,
//			Description:  instance.Description,
//			Contact:      instance.Contact,
//			Type:         instance.Type,
//			Status:       instance.Status,
//			OrgUUID:      instance.OrgUUID,
//			TeamUUID:     instance.TeamUUID,
//			Path:         webFilePath,
//			Icon:         iconPath,
//			Permission:   permissions,
//		}
//		var item = &plugin.PluginConfig{
//			Apis:        apis,
//			Service:     service,
//			Modules:     modules,
//			Abilities:   abilities,
//			PackagePath: instance.PackagePath,
//		}
//		webServiceUrl, _ := appservice.GetLocalPluginWebServiceUrl(instance.OrgUUID, instance.TeamUUID, instance.AppUUID)
//		item.WebServiceUrl = webServiceUrl
//		i18n.HandlePluginWithLanguage(item, req.Language)
//		all = append(all, item)
//	}
//
//	sort.Slice(instances, func(i, j int) bool {
//		return instances[i].AppUUID < instances[j].AppUUID
//	})
//
//	return gin.H{"data": all}, nil
//}
