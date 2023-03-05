package lifecycle

import (
	"github.com/gin-gonic/gin"
	"github.com/obgnail/plugin-platform/common/errors"
	"github.com/obgnail/plugin-platform/common/log"
	"github.com/obgnail/plugin-platform/platform/conn/handler"
	"github.com/obgnail/plugin-platform/platform/model/mysql"
	"github.com/obgnail/plugin-platform/platform/pool/plugin_pool"
	"github.com/obgnail/plugin-platform/platform/service/utils"
	"os"
)

type UninstallReq struct {
	AppUUID      string `json:"app_uuid"`
	InstanceUUID string `json:"instance_uuid"`
}

func (i *UninstallReq) validate() error {
	if i.AppUUID == "" {
		return errors.MissingParameterError(errors.PluginInstanceUninstallationFailure, errors.AppUUID)
	}
	if i.InstanceUUID == "" {
		return errors.MissingParameterError(errors.PluginInstanceUninstallationFailure, errors.InstanceUUID)
	}
	return nil
}

func Uninstall(req *UninstallReq) (ret gin.H, err error) {
	if err := req.validate(); err != nil {
		return ret, errors.Trace(err)
	}
	helper := &UninstallHelper{req: req}
	if err := helper.checkUninstall(); err != nil {
		return ret, errors.Trace(err)
	}
	if err := helper.Uninstall(); err != nil {
		return ret, errors.Trace(err)
	}
	if err := helper.UpdateDb(); err != nil {
		return ret, errors.Trace(err)
	}
	if err := helper.RemoveWorkspace(); err != nil {
		return ret, errors.Trace(err)
	}
	return gin.H{"data": true}, nil
}

type UninstallHelper struct {
	req      *UninstallReq
	pkg      *mysql.PluginPackage
	instance *mysql.PluginInstance
}

func (h *UninstallHelper) checkUninstall() error {
	instanceModel := mysql.ModelPluginInstance()
	instance := &mysql.PluginInstance{AppUUID: h.req.AppUUID, InstanceUUID: h.req.InstanceUUID}
	exist, err := instanceModel.Exist(instance)
	if err != nil {
		log.ErrorDetails(errors.Trace(err))
		return errors.PluginUninstallError(errors.ServerError)
	}
	if exist && instance.Status == plugin_pool.PluginStatusRunning {
		return errors.PluginUninstallError(errors.PluginAlreadyRunning)
	}
	h.instance = instance
	return nil
}

func (h *UninstallHelper) Uninstall() error {
	pck := mysql.ModelPluginPackage()
	pckOne := &mysql.PluginPackage{
		AppUUID: h.req.AppUUID,
		Version: h.instance.Version,
	}
	if err := pck.One(pckOne); err != nil {
		log.ErrorDetails(errors.Trace(err))
		return errors.PluginUninstallError(errors.ServerError)
	}

	h.pkg = pckOne

	cfg, err := pckOne.LoadYamlConfig()
	if err != nil {
		return errors.Trace(err)
	}

	er := handler.UnInstallPlugin(h.req.AppUUID, h.req.InstanceUUID, h.instance.Name,
		cfg.Language, cfg.LanguageVersion, cfg.Version)
	if er != nil {
		log.PEDetails(er)
		return errors.PluginUninstallError(er.Error() + " " + er.Msg())
	}
	return nil
}

func (h *UninstallHelper) UpdateDb() error {
	if err := mysql.ModelPluginPackage().RealDelete(h.pkg.Id, h.pkg); err != nil {
		log.ErrorDetails(errors.Trace(err))
		return errors.PluginUninstallError(errors.ServerError)
	}
	return nil
}

func (h *UninstallHelper) RemoveWorkspace() error {
	path := utils.GetPluginWorkspace(h.instance.AppUUID, h.instance.InstanceUUID)
	if err := os.RemoveAll(path); err != nil {
		log.ErrorDetails(errors.Trace(err))
		return errors.PluginUninstallError(errors.ServerError)
	}
	return nil
}
