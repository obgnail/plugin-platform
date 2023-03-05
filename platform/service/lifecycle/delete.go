package lifecycle

import (
	"github.com/gin-gonic/gin"
	"github.com/obgnail/plugin-platform/common/errors"
	"github.com/obgnail/plugin-platform/platform/model/mysql"
	"github.com/obgnail/plugin-platform/platform/service/utils"
	"os"
)

type DeleteReq struct {
	AppUUID string `json:"app_uuid"`
	Version string `json:"version"`
}

type DeleteResponse struct {
	Result bool `json:"result"`
}

func Delete(req *DeleteReq) (ret gin.H, err error) {
	if err := validate(req); err != nil {
		return ret, errors.Trace(err)
	}
	if err := deletePackage(req); err != nil {
		return ret, errors.Trace(err)
	}
	resp := &DeleteResponse{Result: true}
	return gin.H{"data": resp}, nil
}

func validate(req *DeleteReq) error {
	instanceArg := &mysql.PluginInstance{
		AppUUID: req.AppUUID,
		Version: req.Version,
	}
	exist, err := mysql.ModelPluginInstance().Exist(instanceArg)
	if err != nil {
		return errors.Trace(err)
	}
	// 运行中的实例不能删除
	if exist {
		return errors.Trace(err)
	}
	return nil
}

func deletePackage(req *DeleteReq) error {
	m := mysql.ModelPluginPackage()

	one := &mysql.PluginPackage{
		AppUUID: req.AppUUID,
		Version: req.Version,
	}

	if err := m.One(one); err != nil {
		return errors.Trace(err)
	}

	if err := m.RealDelete(one.Id, one); err != nil {
		return errors.Trace(err)
	}
	path := utils.GetPluginDir(one.AppUUID, one.Version)
	if err := os.RemoveAll(path); err != nil {
		return errors.Trace(err)
	}
	return nil
}
