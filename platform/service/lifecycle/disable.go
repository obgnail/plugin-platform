package lifecycle

import (
	"github.com/gin-gonic/gin"
	"github.com/obgnail/plugin-platform/common/errors"
	"github.com/obgnail/plugin-platform/common/log"
	"github.com/obgnail/plugin-platform/platform/conn/handler"
	"github.com/obgnail/plugin-platform/platform/model/mysql"
	"github.com/obgnail/plugin-platform/platform/service/common"
)

type DisableReq struct {
	InstanceUUID string `json:"instance_uuid"`
}

func (i *DisableReq) validate() error {
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
	instance := &mysql.PluginInstance{InstanceUUID: h.req.InstanceUUID}
	exist, err := mysql.ModelPluginInstance().Exist(instance)
	if err != nil {
		log.ErrorDetails(errors.Trace(err))
		return errors.PluginDisableError(errors.ServerError)
	}
	if exist && instance.Status != common.PluginStatusRunning {
		return errors.PluginDisableError(errors.PluginAlreadyStop)
	}
	h.instance = instance
	return nil
}

func (h *DisableHelper) Disable() error {
	cnf, err := h.instance.LoadYamlConfig()
	if err != nil {
		return errors.Trace(err)
	}

	er := handler.DisablePlugin(h.instance.AppUUID, h.instance.InstanceUUID, h.instance.Name,
		cnf.Language, cnf.LanguageVersion, cnf.Version)
	if er != nil {
		log.PEDetails(er)
		return errors.PluginDisableError(er.Error() + " " + er.Msg())
	}
	return nil
}

func (h *DisableHelper) UpdateDb() error {
	h.instance.Status = common.PluginStatusStopping
	if err := mysql.ModelPluginInstance().Update(h.instance.Id, h.instance); err != nil {
		log.ErrorDetails(errors.Trace(err))
		return errors.PluginDisableError(errors.ServerError)
	}
	return nil
}
