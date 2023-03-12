package lifecycle

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/obgnail/plugin-platform/common/errors"
	"github.com/obgnail/plugin-platform/common/log"
	"github.com/obgnail/plugin-platform/common/utils/math"
	"github.com/obgnail/plugin-platform/platform/conn/handler"
	"github.com/obgnail/plugin-platform/platform/model/mysql"
	"github.com/obgnail/plugin-platform/platform/service/common"
)

type InstallReq struct {
	InstanceUUID string `json:"instance_uuid"`
}

type InstallResp struct {
	*common.Service `json:"service"`
	Apis            []*common.Api     `json:"apis"`
	Abilities       []*common.Ability `json:"abilities"`
}

func (i *InstallReq) validate() error {
	if i.InstanceUUID == "" {
		return errors.MissingParameterError(errors.PluginInstanceInstallationFailure, errors.InstanceUUID)
	}
	return nil
}

func Install(req *InstallReq) (ret gin.H, err error) {
	if err = req.validate(); err != nil {
		return ret, errors.Trace(err)
	}
	helper := &InstallHelper{req: req}
	if err := helper.checkInstall(); err != nil {
		return ret, errors.Trace(err)
	}
	if err := helper.generatePlugin(); err != nil {
		return ret, errors.Trace(err)
	}
	if err := helper.Save2Db(); err != nil {
		return ret, errors.Trace(err)
	}
	var resp = &InstallResp{
		Service:   helper.Cfg.Service,
		Apis:      helper.Cfg.Apis,
		Abilities: helper.Cfg.Abilities,
	}
	return gin.H{"data": resp}, err
}

type InstallHelper struct {
	req      *InstallReq
	Cfg      *common.PluginConfig
	instance *mysql.PluginInstance
}

func (h *InstallHelper) checkInstall() error {
	instanceModel := mysql.ModelPluginInstance()
	instance := &mysql.PluginInstance{InstanceUUID: h.req.InstanceUUID}
	exist, err := instanceModel.Exist(instance)
	if err != nil {
		log.ErrorDetails(errors.Trace(err))
		return errors.PluginInstallError(errors.ServerError)
	}
	if exist && instance.Status != common.PluginStatusUploaded {
		return errors.PluginInstallError(errors.PluginAlreadyInstall)
	}
	h.instance = instance
	return nil
}

func (h *InstallHelper) generatePlugin() (err error) {
	h.Cfg, err = h.instance.LoadYamlConfig()
	if err != nil {
		return errors.PluginInstallError(errors.LoadYamlConfigFailed)
	}
	h.Cfg.Service.InstanceUUID = h.req.InstanceUUID
	er := <-handler.InstallPlugin(h.instance.AppUUID, h.instance.InstanceUUID, h.instance.Name,
		h.Cfg.Language, h.Cfg.LanguageVersion, h.Cfg.Version)
	if er != nil {
		log.PEDetails(er)
		return errors.PluginInstallError(er.Error() + " " + er.Msg())
	}
	return nil
}

func (h *InstallHelper) Save2Db() error {
	err := mysql.Transaction(func(db *gorm.DB) error {
		h.instance.Status = common.PluginStatusStopping
		if err := mysql.ModelPluginInstance().Update(h.instance.Id, h.instance); err != nil {
			return errors.Trace(err)
		}

		// 生成插件配置
		if err := generatePluginConfig(db, h.instance, h.Cfg); err != nil {
			return errors.Trace(err)
		}

		// 生成自定义权限点
		if err := generatePluginPermission(db, h.instance, h.Cfg); err != nil {
			return errors.Trace(err)
		}

		// 给每个插件生成一个新的用户,方便标品鉴权
		if err := generatePluginUser(db, h.instance); err != nil {
			return errors.Trace(err)
		}
		return nil
	})
	if err != nil {
		log.ErrorDetails(errors.Trace(err))
		return errors.PluginInstallError(errors.SaveToDBFailed)
	}
	return nil
}

func generatePluginConfig(db *gorm.DB, instance *mysql.PluginInstance, cfg *common.PluginConfig) error {
	configs := cfg.Service.Config
	if len(configs) == 0 {
		return nil
	}

	var dataset = make([]*mysql.PluginConfig, 0)
	for _, c := range configs {
		var d = &mysql.PluginConfig{
			AppUUID:      instance.AppUUID,
			InstanceUUID: instance.InstanceUUID,
			Label:        c.Label,
			Key:          c.Key,
			Value:        c.Value,
			Type:         mysql.ConvertConfigType(c.Type),
			Required:     c.Required,
		}
		dataset = append(dataset, d)
	}

	if err := mysql.ModelPluginConfig().NewBatchWithDB(db, dataset); err != nil {
		return errors.Trace(err)
	}
	return nil
}

func generatePluginPermission(db *gorm.DB, instance *mysql.PluginInstance, cfg *common.PluginConfig) error {
	permission := cfg.Service.Permission
	if len(permission) == 0 {
		return nil
	}
	var permissionData = make([]*mysql.PluginPermissionInfo, 0)
	for _, info := range permission {
		m := &mysql.PluginPermissionInfo{
			InstanceUUID:    instance.InstanceUUID,
			PermissionName:  info.Name,
			PermissionField: info.Field,
			PermissionDesc:  info.Desc,
			PermissionID:    int(math.CreateCaptcha()),
		}
		permissionData = append(permissionData, m)
	}

	if err := mysql.ModelPluginPermissionInfo().NewBatchWithDB(db, permissionData); err != nil {
		return errors.Trace(err)
	}
	return nil
}

func generatePluginUser(db *gorm.DB, instance *mysql.PluginInstance) error {
	u := &mysql.PluginUser{
		UserUUID:     mysql.NewUserUUID(instance.AppUUID, instance.InstanceUUID),
		AppUUID:      instance.AppUUID,
		InstanceUUID: instance.InstanceUUID,
		Name:         instance.Name,
		Email:        mysql.NewUserEmail(instance.AppUUID, instance.InstanceUUID),
	}
	dataset := []*mysql.PluginUser{u}
	if err := mysql.ModelPluginUser().NewBatchWithDB(db, dataset); err != nil {
		return errors.Trace(err)
	}
	return nil
}
