package message

import (
	"github.com/gin-gonic/gin"
	"github.com/obgnail/plugin-platform/common/errors"
	"github.com/obgnail/plugin-platform/common/log"
	"github.com/obgnail/plugin-platform/platform/model/mysql"
	"github.com/obgnail/plugin-platform/platform/service/types"
)

type Plugin struct {
	UUID        string           `json:"uuid"`
	Name        string           `json:"name"`
	Version     string           `json:"version"`
	LifeStage   int              `json:"life_stage"`
	Description string           `json:"description"`
	Routers     []*types.Api     `json:"routers"`
	Abilities   []*types.Ability `json:"abilities"`
}

func RouterList() (ret gin.H, err error) {
	var instances = make([]*mysql.PluginInstance, 0)
	err = mysql.ModelPluginInstance().All(&instances, &mysql.PluginInstance{})
	if err != nil {
		log.ErrorDetails(errors.Trace(err))
		return ret, errors.PluginMessageError(errors.ServerError)
	}

	var plugins []*Plugin
	for _, instance := range instances {
		cfg, err := instance.LoadYamlConfig()
		if err != nil {
			log.ErrorDetails(errors.Trace(err))
			return ret, errors.PluginMessageError(errors.ServerError)
		}

		plugin := &Plugin{
			UUID:        instance.InstanceUUID,
			Name:        instance.Name,
			Version:     instance.Version,
			LifeStage:   instance.Status,
			Description: instance.Description,
			Routers:     cfg.Apis,
			Abilities:   cfg.Abilities,
		}
		plugins = append(plugins, plugin)
	}
	return gin.H{"data": plugins}, nil
}
