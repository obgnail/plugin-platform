package lifecycle

import (
	"github.com/gin-gonic/gin"
	"github.com/obgnail/plugin-platform/common/errors"
	"github.com/obgnail/plugin-platform/common/log"
	"github.com/obgnail/plugin-platform/platform/model/mysql"
	"github.com/obgnail/plugin-platform/platform/pool/plugin_pool"
	"github.com/obgnail/plugin-platform/platform/service/types"
	"github.com/obgnail/plugin-platform/platform/service/utils"
	"gopkg.in/yaml.v2"
	"mime/multipart"
	"path/filepath"
)

type UpgradeReq struct {
	AppUUID    string `json:"app_uuid"`
	AppVersion string `json:"app_version"`
}

func (i *UpgradeReq) validate() error {
	if i.AppUUID == "" {
		return errors.MissingParameterError(errors.PluginInstanceUpgradeFailure, errors.AppUUID)
	}
	if i.AppVersion == "" {
		return errors.MissingParameterError(errors.PluginInstanceUpgradeFailure, errors.AppVersion)
	}
	return nil
}

// Upgrade 升级操作 不重新生成instanceID, 对插件包进行替换操作
func Upgrade(ctx *gin.Context, req *UpgradeReq) (ret gin.H, err error) {
	if err := req.validate(); err != nil {
		return ret, errors.Trace(err)
	}
	helper := &UpgradeHelper{ctx: ctx, req: req}
	if err = helper.checkFile(); err != nil {
		return ret, errors.Trace(err)
	}
	if err = helper.parseConfigYaml(); err != nil {
		return ret, errors.Trace(err)
	}
	if err = helper.checkUpgrade(); err != nil {
		return ret, errors.Trace(err)
	}
	if err = helper.upgrade(); err != nil {
		return ret, errors.Trace(err)
	}
	return gin.H{"data": true}, nil
}

type UpgradeHelper struct {
	ctx *gin.Context
	req *UpgradeReq

	fileHeader *multipart.FileHeader
	file       multipart.File

	newConfig *plugin_pool.PluginConfig // config.yaml

	oldInstance *mysql.PluginInstance
	oldPackage  *mysql.PluginPackage
	oldConfig   *plugin_pool.PluginConfig
}

func (h *UpgradeHelper) checkFile() error {
	err := CheckMaxFileSize(h.ctx, types.PluginFileMaxSize)
	if err != nil {
		return errors.PluginUpgradeError(errors.FileTooLarge)
	}
	h.file, h.fileHeader, err = h.ctx.Request.FormFile("file")
	if err != nil {
		return errors.PluginUpgradeError(errors.FileUploadFailed)
	}
	if filepath.Ext(h.fileHeader.Filename) != types.PluginFileExt {
		return errors.PluginUpgradeError(errors.InvalidFileExt)
	}
	return nil
}

func (h *UpgradeHelper) parseConfigYaml() error {
	configYaml, err := utils.GetYamlFromFile(h.fileHeader, h.file)
	if err != nil {
		return errors.PluginUpgradeError(errors.PluginConfigFileParseFailed)
	}

	var pluginConfig = new(plugin_pool.PluginConfig)
	if err = yaml.Unmarshal([]byte(configYaml), pluginConfig); err != nil {
		log.Warn("Unmarshal err: %+v", errors.Trace(err))
		return errors.PluginUpgradeError(errors.PluginConfigFileParseFailed)
	}

	h.newConfig = pluginConfig
	return nil
}
func (h *UpgradeHelper) checkUpgrade() error {
	// 判断AppUUID跟AppVersion是否还是一致,新的插件有自己新的AppUUID和AppVersion
	if h.req.AppUUID == h.newConfig.AppUUID || h.req.AppVersion == h.newConfig.Version {
		return errors.PluginUpgradeError(errors.ErrorMetaData)
	}

	instanceModel := mysql.ModelPluginInstance()
	oldInstance := &mysql.PluginInstance{
		AppUUID: h.req.AppUUID,
		Version: h.req.AppVersion,
	}
	exist1, err := instanceModel.Exist(oldInstance)
	if err != nil {
		log.ErrorDetails(errors.Trace(err))
		return errors.PluginUpgradeError(errors.ServerError)
	}
	if !exist1 {
		return errors.PluginUpgradeError(errors.InstanceNotFound)
	}
	h.oldInstance = oldInstance

	pkgModel := mysql.ModelPluginPackage()
	oldPackage := &mysql.PluginPackage{
		AppUUID: h.req.AppUUID,
		Version: h.req.AppVersion,
	}
	exist2, err := pkgModel.Exist(oldPackage)
	if err != nil {
		log.ErrorDetails(errors.Trace(err))
		return errors.PluginUpgradeError(errors.ServerError)
	}
	if !exist2 {
		return errors.PluginUpgradeError(errors.PackageNotFound)
	}
	h.oldPackage = oldPackage

	config, err := oldPackage.LoadYamlConfig()
	if err != nil {
		return errors.Trace(err)
	}
	h.oldConfig = config
	return nil
}

func (h *UpgradeHelper) upgrade() error {
	//er := handler.UpgradePlugin(h.oldInstance.AppUUID, h.oldInstance.InstanceUUID, h.oldInstance.Name,
	//	h.oldConfig.Language, h.oldConfig.Version)
	//if er != nil {
	//	log.PEDetails(er)
	//	return errors.PluginDisableError(er.Error() + " " + er.Msg())
	//}

	return nil
}
