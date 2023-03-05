package lifecycle

import (
	"github.com/gin-gonic/gin"
	"github.com/obgnail/plugin-platform/common/errors"
	"github.com/obgnail/plugin-platform/common/log"
	"github.com/obgnail/plugin-platform/platform/conn/handler"
	"github.com/obgnail/plugin-platform/platform/model/mysql"
	"github.com/obgnail/plugin-platform/platform/pool/plugin_pool"
)

type EnableReq struct {
	AppUUID      string `json:"app_uuid"`
	InstanceUUID string `json:"instance_uuid"`
}

func (i *EnableReq) validate() error {
	if i.AppUUID == "" {
		return errors.MissingParameterError(errors.PluginInstanceEnableFailure, errors.AppUUID)
	}
	if i.InstanceUUID == "" {
		return errors.MissingParameterError(errors.PluginInstanceEnableFailure, errors.InstanceUUID)
	}
	return nil
}

type EnableHelper struct {
	req      *EnableReq
	instance *mysql.PluginInstance
}

func Enable(req *EnableReq) (ret gin.H, err error) {
	if err := req.validate(); err != nil {
		return ret, errors.Trace(err)
	}

	helper := &EnableHelper{req: req}
	if err := helper.checkEnable(); err != nil {
		return ret, errors.Trace(err)
	}
	if err := helper.Enable(); err != nil {
		return ret, errors.Trace(err)
	}
	if err := helper.UpdateDb(); err != nil {
		return ret, errors.Trace(err)
	}
	return gin.H{"data": true}, nil
}

func (h *EnableHelper) checkEnable() error {
	instanceModel := mysql.ModelPluginInstance()
	instance := &mysql.PluginInstance{AppUUID: h.req.AppUUID, InstanceUUID: h.req.InstanceUUID}
	exist, err := instanceModel.Exist(instance)
	if err != nil {
		log.ErrorDetails(errors.Trace(err))
		return errors.PluginEnableError(errors.ServerError)
	}
	if exist && instance.Status == plugin_pool.PluginStatusRunning {
		return errors.PluginEnableError(errors.PluginAlreadyRunning)
	}
	h.instance = instance
	return nil
}

func (h *EnableHelper) Enable() error {
	pck := mysql.ModelPluginPackage()
	pckOne := &mysql.PluginPackage{
		AppUUID: h.req.AppUUID,
		Version: h.instance.Version,
	}
	if err := pck.One(pckOne); err != nil {
		log.ErrorDetails(errors.Trace(err))
		return errors.PluginEnableError(errors.ServerError)
	}

	config, err := pckOne.LoadYamlConfig()
	if err != nil {
		return errors.Trace(err)
	}

	er := handler.EnablePlugin(h.req.AppUUID, h.req.InstanceUUID, h.instance.Name,
		config.Language, config.LanguageVersion, config.Version)
	if er != nil {
		log.PEDetails(er)
		return errors.PluginEnableError(er.Error() + " " + er.Msg())
	}
	return nil
}

func (h *EnableHelper) UpdateDb() error {
	h.instance.Status = plugin_pool.PluginStatusRunning
	if err := mysql.ModelPluginInstance().Update(h.instance.Id, h.instance); err != nil {
		log.ErrorDetails(errors.Trace(err))
		return errors.PluginEnableError(errors.ServerError)
	}
	return nil
}
