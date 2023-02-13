package lifecycle

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/obgnail/plugin-platform/platform/model/mysql"
	"github.com/obgnail/plugin-platform/utils/errors"
	"github.com/obgnail/plugin-platform/utils/log"
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
		return
	}
	helper := &InstallHelper{req: req}
	pkg, err := helper.checkInstall()
	if err != nil {
		return ret, errors.Trace(err)
	}
	if err := helper.generatePlugin(pkg); err != nil {
		return ret, errors.Trace(err)
	}

	return gin.H{"data": "resp"}, err
}

type InstallHelper struct {
	req *InstallReq
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

func (h *InstallHelper) generatePlugin(pkg *mysql.PluginPackage) error {
	cfg, err := pkg.LoadYamlConfig()
	if err != nil {
		return errors.Trace(err)
	}
	fmt.Println(cfg)
	return nil
}
