package lifecycle

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/obgnail/plugin-platform/common/errors"
	"github.com/obgnail/plugin-platform/common/log"
	"github.com/obgnail/plugin-platform/common/utils/math"
	"github.com/obgnail/plugin-platform/platform/conn/handler"
	"github.com/obgnail/plugin-platform/platform/model/mysql"
	"github.com/obgnail/plugin-platform/platform/pool/plugin_pool"
	"github.com/obgnail/plugin-platform/platform/service/utils"
)

type InstallReq struct {
	AppUUID      string `json:"app_uuid"`
	InstanceUUID string `json:"instance_uuid"`
}

func (i *InstallReq) validate() error {
	if i.AppUUID == "" {
		return errors.MissingParameterError(errors.PluginInstanceInstallationFailure, errors.AppUUID)
	}
	return nil
}

func Install(req *InstallReq) (ret gin.H, err error) {
	if err = req.validate(); err != nil {
		return ret, errors.Trace(err)
	}
	helper := &InstallHelper{req: req}
	pkg, err := helper.checkInstall()
	if err != nil {
		return ret, errors.Trace(err)
	}
	if err := helper.generatePlugin(pkg); err != nil {
		return ret, errors.Trace(err)
	}
	if err := helper.Save2Db(); err != nil {
		return ret, errors.Trace(err)
	}

	return gin.H{"data": "resp"}, err
}

type InstallHelper struct {
	req          *InstallReq
	cfg          *plugin_pool.PluginConfig
	instanceUUID string
}

func (h *InstallHelper) checkInstall() (*mysql.PluginPackage, error) {
	pkgModel := mysql.ModelPluginPackage()
	pkg := &mysql.PluginPackage{AppUUID: h.req.AppUUID}
	err := pkgModel.One(pkg)
	if err == mysql.RecordNotFound {
		return nil, errors.PluginInstallError(errors.FileNoExist)
	}
	if err != nil {
		log.ErrorDetails(errors.Trace(err))
		return nil, errors.PluginInstallError(errors.ServerError)
	}

	instanceModel := mysql.ModelPluginInstance()
	instance := &mysql.PluginInstance{AppUUID: h.req.AppUUID}
	exist, err := instanceModel.Exist(instance)
	if err != nil {
		log.ErrorDetails(errors.Trace(err))
		return nil, errors.PluginInstallError(errors.ServerError)
	}
	if exist {
		return nil, errors.PluginInstallError(errors.PluginAlreadyInstall)
	}
	return pkg, nil
}

func (h *InstallHelper) generatePlugin(pkg *mysql.PluginPackage) (err error) {
	h.cfg, err = pkg.LoadYamlConfig()
	if err != nil {
		return errors.PluginInstallError(errors.LoadYamlConfigFailed)
	}
	h.cfg.Status = plugin_pool.PluginStatusRunning
	h.instanceUUID = utils.NewInstanceUUID()
	er := handler.InstallPlugin(h.cfg.AppUUID, h.instanceUUID, h.cfg.Name,
		h.cfg.Language, h.cfg.LanguageVersion, h.cfg.Version)
	if er != nil {
		err = errors.PluginInstallError(er.Error())
	}
	return nil
}

func (h *InstallHelper) Save2Db() error {
	apis, err := json.Marshal(h.cfg.Apis)
	if err != nil {
		log.ErrorDetails(errors.Trace(err))
		return errors.PluginInstallError(errors.LoadYamlConfigFailed)
	}

	m := &mysql.PluginInstance{
		AppUUID:      h.cfg.AppUUID,
		InstanceUUID: h.cfg.InstanceUUID,
		Name:         h.cfg.Name,
		Version:      h.cfg.Version,
		Description:  h.cfg.Description,
		Contact:      h.cfg.Contact,
		Status:       h.cfg.Status,
		Apis:         string(apis),
	}
	models := []*mysql.PluginInstance{m}

	err = mysql.Transaction(func(db *gorm.DB) error {
		if err = mysql.ModelPluginInstance().NewBatchWithDB(db, models); err != nil {
			return errors.Trace(err)
		}

		// 生成插件配置
		if err = generatePluginConfig(db, models, h.cfg); err != nil {
			return errors.Trace(err)
		}

		// 生成自定义权限点
		if err = generatePluginPermission(db, h.cfg); err != nil {
			return errors.Trace(err)
		}

		// 给每个插件生成一个新的用户,方便标品鉴权
		if err = generatePluginUser(db, h.cfg); err != nil {
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

func generatePluginConfig(db *gorm.DB, instances []*mysql.PluginInstance, cfg *plugin_pool.PluginConfig) error {
	configs := cfg.Service.Config
	if len(configs) == 0 {
		return nil
	}

	var dataset = make([]*mysql.PluginConfig, 0)
	for _, c := range configs {
		for _, instance := range instances {
			var d = &mysql.PluginConfig{
				AppUUID:      instance.AppUUID,
				InstanceUUID: instance.InstanceUUID,
				Key:          c.Key,
				Value:        c.Value,
				Type:         mysql.ConvertConfigType(c.Type),
				Required:     c.Required,
			}
			dataset = append(dataset, d)
		}
	}

	if err := mysql.ModelPluginConfig().NewBatchWithDB(db, dataset); err != nil {
		return errors.Trace(err)
	}
	return nil
}

func generatePluginPermission(db *gorm.DB, cfg *plugin_pool.PluginConfig) error {
	permission := cfg.Service.Permission
	if len(permission) == 0 {
		return nil
	}
	var permissionData = make([]*mysql.PluginPermissionInfo, 0)
	for _, info := range permission {
		m := &mysql.PluginPermissionInfo{
			InstanceUUID:    cfg.InstanceUUID,
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

func generatePluginUser(db *gorm.DB, cfg *plugin_pool.PluginConfig) error {
	u := &mysql.PluginUser{
		UserUUID:     mysql.NewUserUUID(cfg.AppUUID, cfg.InstanceUUID),
		AppUUID:      cfg.AppUUID,
		InstanceUUID: cfg.InstanceUUID,
		Name:         cfg.Name,
		Email:        mysql.NewUserEmail(cfg.AppUUID, cfg.InstanceUUID),
	}
	dataset := []*mysql.PluginUser{u}
	if err := mysql.ModelPluginConfig().NewBatchWithDB(db, dataset); err != nil {
		return errors.Trace(err)
	}
	return nil
}
