package lifecycle

import (
	"github.com/gin-gonic/gin"
	"github.com/obgnail/plugin-platform/common/errors"
	"github.com/obgnail/plugin-platform/common/log"
	"github.com/obgnail/plugin-platform/platform/conn/ability"
	"github.com/obgnail/plugin-platform/platform/conn/handler"
	"github.com/obgnail/plugin-platform/platform/conn/router"
	"github.com/obgnail/plugin-platform/platform/model/mysql"
	"github.com/obgnail/plugin-platform/platform/service/common"
)

type EnableReq struct {
	InstanceUUID string `json:"instance_uuid"`
}

func (i *EnableReq) validate() error {
	if i.InstanceUUID == "" {
		return errors.MissingParameterError(errors.PluginInstanceEnableFailure, errors.InstanceUUID)
	}
	return nil
}

type EnableHelper struct {
	req      *EnableReq
	cfg      *common.PluginConfig
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
	if err := helper.RegisterRouter(); err != nil {
		return ret, errors.Trace(err)
	}
	if err := helper.RegisterAbility(); err != nil {
		return ret, errors.Trace(err)
	}
	if err := helper.UpdateDb(); err != nil {
		return ret, errors.Trace(err)
	}
	return gin.H{"data": true}, nil
}

func (h *EnableHelper) checkEnable() error {
	instanceModel := mysql.ModelPluginInstance()
	instance := &mysql.PluginInstance{InstanceUUID: h.req.InstanceUUID}
	exist, err := instanceModel.Exist(instance)
	if err != nil {
		log.ErrorDetails(errors.Trace(err))
		return errors.PluginEnableError(errors.ServerError)
	}
	if exist && instance.Status != common.PluginStatusStopping {
		return errors.PluginEnableError(errors.PluginAlreadyRunning)
	}
	h.instance = instance
	return nil
}

func (h *EnableHelper) Enable() error {
	cfg, err := h.instance.LoadYamlConfig()
	if err != nil {
		return errors.Trace(err)
	}

	er := <-handler.EnablePlugin(h.instance.InstanceUUID)
	if er != nil {
		log.PEDetails(er)
		return errors.PluginEnableError(er.Error() + " " + er.Msg())
	}
	h.cfg = cfg
	return nil
}

func (h *EnableHelper) RegisterRouter() error {
	apis := h.cfg.Apis
	if err := router.RegisterRouter(apis, h.instance.InstanceUUID); err != nil {
		log.ErrorDetails(errors.Trace(err))
		return errors.PluginEnableError(errors.ServerError)
	}
	return nil
}

func (h *EnableHelper) RegisterAbility() error {
	abilities := h.cfg.Abilities
	if err := ability.RegisterAbility(abilities, h.instance.InstanceUUID); err != nil {
		log.ErrorDetails(errors.Trace(err))
		return errors.PluginEnableError(errors.ServerError)
	}
	return nil
}

func (h *EnableHelper) UpdateDb() error {
	h.instance.Status = common.PluginStatusRunning
	if err := mysql.ModelPluginInstance().Update(h.instance.Id, h.instance); err != nil {
		log.ErrorDetails(errors.Trace(err))
		return errors.PluginEnableError(errors.ServerError)
	}
	return nil
}
