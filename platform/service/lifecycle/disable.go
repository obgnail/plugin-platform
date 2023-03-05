package lifecycle

import (
	"github.com/gin-gonic/gin"
	"github.com/obgnail/plugin-platform/common/errors"
	"github.com/obgnail/plugin-platform/common/log"
	"github.com/obgnail/plugin-platform/platform/conn/handler"
	"github.com/obgnail/plugin-platform/platform/model/mysql"
	"github.com/obgnail/plugin-platform/platform/pool/plugin_pool"
)

type DisableReq struct {
	AppUUID      string `json:"app_uuid"`
	InstanceUUID string `json:"instance_uuid"`
}

func (i *DisableReq) validate() error {
	if i.AppUUID == "" {
		return errors.MissingParameterError(errors.PluginInstanceDisableFailure, errors.AppUUID)
	}
	if i.InstanceUUID == "" {
		return errors.MissingParameterError(errors.PluginInstanceDisableFailure, errors.InstanceUUID)
	}
	return nil
}

type DisableHelper struct {
	req      *DisableReq
	instance *mysql.PluginInstance
}

func Disable(req *DisableReq) (ret gin.H, err error) {
	if err := req.validate(); err != nil {
		return ret, errors.Trace(err)
	}
	helper := &DisableHelper{req: req}
	if err := helper.checkDisable(); err != nil {
		return ret, errors.Trace(err)
	}
	if err := helper.Disable(); err != nil {
		return ret, errors.Trace(err)
	}
	if err := helper.UpdateDb(); err != nil {
		return ret, errors.Trace(err)
	}
	return gin.H{"data": true}, nil
}

func (h *DisableHelper) checkDisable() error {
	instanceModel := mysql.ModelPluginInstance()
	instance := &mysql.PluginInstance{AppUUID: h.req.AppUUID, InstanceUUID: h.req.InstanceUUID}
	exist, err := instanceModel.Exist(instance)
	if err != nil {
		log.ErrorDetails(errors.Trace(err))
		return errors.PluginDisableError(errors.ServerError)
	}
	if exist && instance.Status == plugin_pool.PluginStatusStopping {
		return errors.PluginDisableError(errors.PluginAlreadyStop)
	}
	h.instance = instance
	return nil
}

func (h *DisableHelper) Disable() error {
	pck := mysql.ModelPluginPackage()
	pckOne := &mysql.PluginPackage{
		AppUUID: h.req.AppUUID,
		Version: h.instance.Version,
	}
	if err := pck.One(pckOne); err != nil {
		log.ErrorDetails(errors.Trace(err))
		return errors.PluginDisableError(errors.ServerError)
	}

	config, err := pckOne.LoadYamlConfig()
	if err != nil {
		return errors.Trace(err)
	}

	er := handler.DisablePlugin(h.req.AppUUID, h.req.InstanceUUID, h.instance.Name,
		config.Language, config.LanguageVersion, config.Version)
	if er != nil {
		log.PEDetails(er)
		return errors.PluginDisableError(er.Error() + " " + er.Msg())
	}
	return nil
}

func (h *DisableHelper) UpdateDb() error {
	h.instance.Status = plugin_pool.PluginStatusStopping
	if err := mysql.ModelPluginInstance().Update(h.instance.Id, h.instance); err != nil {
		log.ErrorDetails(errors.Trace(err))
		return errors.PluginDisableError(errors.ServerError)
	}
	return nil
}
