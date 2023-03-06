package message

import (
	"github.com/gin-gonic/gin"
	"github.com/obgnail/plugin-platform/common/errors"
	"github.com/obgnail/plugin-platform/common/log"
	"github.com/obgnail/plugin-platform/platform/model/mysql"
	"github.com/obgnail/plugin-platform/platform/service/common"
)

type Item struct {
	AppUUID      string               `json:"app_uuid"`
	InstanceUUID string               `json:"instance_uuid"`
	Name         string               `json:"name"`
	Version      string               `json:"version"`
	Description  string               `json:"description"`
	Status       int                  `json:"status"`
	Permission   []*common.Permission `json:"permission"`
}

func ListPlugins() (ret gin.H, err error) {
	var instances = make([]*mysql.PluginInstance, 0)
	err = mysql.ModelPluginInstance().All(&instances, &mysql.PluginInstance{})
	if err != nil {
		log.ErrorDetails(errors.Trace(err))
		return ret, errors.PluginMessageError(errors.ServerError)
	}

	var items []*Item
	for _, instance := range instances {
		var infos = make([]*mysql.PluginPermissionInfo, 0)
		err = mysql.ModelPluginPermissionInfo().All(&infos, &mysql.PluginPermissionInfo{InstanceUUID: instance.InstanceUUID})
		if err != nil {
			log.ErrorDetails(errors.Trace(err))
			return ret, errors.PluginMessageError(errors.ServerError)
		}

		var permission []*common.Permission
		for _, info := range infos {
			permission = append(permission, &common.Permission{
				Name:  info.PermissionName,
				Field: info.PermissionField,
				Desc:  info.PermissionDesc,
			})
		}

		items = append(items, &Item{
			AppUUID:      instance.AppUUID,
			InstanceUUID: instance.InstanceUUID,
			Name:         instance.Name,
			Version:      instance.Version,
			Description:  instance.Description,
			Status:       instance.Status,
			Permission:   permission,
		})
	}

	return gin.H{"data": items}, nil
}
