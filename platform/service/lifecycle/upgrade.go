package lifecycle

//import (
//	"github.com/gin-gonic/gin"
//	"github.com/obgnail/plugin-platform/common/errors"
//	"github.com/obgnail/plugin-platform/common/log"
//	"github.com/obgnail/plugin-platform/common/protocol"
//	"github.com/obgnail/plugin-platform/common/utils/message"
//	"github.com/obgnail/plugin-platform/platform/conn/handler"
//	"github.com/obgnail/plugin-platform/platform/model/mysql"
//	"github.com/obgnail/plugin-platform/platform/pool/plugin_pool"
//	"github.com/obgnail/plugin-platform/platform/service/types"
//	"github.com/obgnail/plugin-platform/platform/service/utils"
//	"gopkg.in/yaml.v2"
//	"mime/multipart"
//	"path/filepath"
//)
//
//type UpgradeReq struct {
//	OldAppUUID    string `json:"app_uuid"`
//	OldAppVersion string `json:"app_version"`
//}
//
//func (i *UpgradeReq) validate() error {
//	if i.OldAppUUID == "" {
//		return errors.MissingParameterError(errors.PluginInstanceUpgradeFailure, errors.AppUUID)
//	}
//	if i.OldAppVersion == "" {
//		return errors.MissingParameterError(errors.PluginInstanceUpgradeFailure, errors.AppVersion)
//	}
//	return nil
//}
//
//// Upgrade 升级操作 不重新生成instanceID, 对插件包进行替换操作
////插件的升级流程：
////	1. 停用旧的插件。并把该插件置为不可用状态
////	2. 上传新的插件包。
////	3. 解压新的插件包。
////  4. 执行新的插件的 Upgrade 生命周期函数。
////  5. 修改plugin_instance、plugin_package、plugin_permission,plugin_config表中相关数据
//func Upgrade(ctx *gin.Context, req *UpgradeReq) (ret gin.H, err error) {
//	if err := req.validate(); err != nil {
//		return ret, errors.Trace(err)
//	}
//	helper := &UpgradeHelper{ctx: ctx, req: req}
//	if err = helper.checkFile(); err != nil {
//		return ret, errors.Trace(err)
//	}
//	if err = helper.parseConfigYaml(); err != nil {
//		return ret, errors.Trace(err)
//	}
//	if err = helper.checkUpgrade(); err != nil {
//		return ret, errors.Trace(err)
//	}
//	if err = helper.upgrade(); err != nil {
//		return ret, errors.Trace(err)
//	}
//	return gin.H{"data": true}, nil
//}
//
//type UpgradeHelper struct {
//	ctx *gin.Context
//	req *UpgradeReq
//
//	fileHeader *multipart.FileHeader
//	file       multipart.File
//
//	newConfig *plugin_pool.PluginConfig // config.yaml
//
//	oldInstance *mysql.PluginInstance
//	oldPackage  *mysql.PluginPackage
//	oldConfig   *plugin_pool.PluginConfig
//}
//
//func (h *UpgradeHelper) checkFile() error {
//	err := CheckMaxFileSize(h.ctx, types.PluginFileMaxSize)
//	if err != nil {
//		return errors.PluginUpgradeError(errors.FileTooLarge)
//	}
//	h.file, h.fileHeader, err = h.ctx.Request.FormFile("file")
//	if err != nil {
//		return errors.PluginUpgradeError(errors.FileUploadFailed)
//	}
//	if filepath.Ext(h.fileHeader.Filename) != types.PluginFileExt {
//		return errors.PluginUpgradeError(errors.InvalidFileExt)
//	}
//	return nil
//}
//
//func (h *UpgradeHelper) parseConfigYaml() error {
//	configYaml, err := utils.GetYamlFromFile(h.fileHeader, h.file)
//	if err != nil {
//		return errors.PluginUpgradeError(errors.PluginConfigFileParseFailed)
//	}
//
//	var pluginConfig = new(plugin_pool.PluginConfig)
//	if err = yaml.Unmarshal([]byte(configYaml), pluginConfig); err != nil {
//		log.Warn("Unmarshal err: %+v", errors.Trace(err))
//		return errors.PluginUpgradeError(errors.PluginConfigFileParseFailed)
//	}
//
//	h.newConfig = pluginConfig
//	return nil
//}
//
//func (h *UpgradeHelper) checkUpgrade() error {
//	// 判断AppUUID跟AppVersion是否还是一致,新的插件有自己新的AppUUID和AppVersion
//	if h.req.OldAppUUID == h.newConfig.AppUUID || h.req.OldAppVersion == h.newConfig.Version {
//		return errors.PluginUpgradeError(errors.ErrorMetaData)
//	}
//
//	// 判断oldPackage是否存在
//	pkgModel := mysql.ModelPluginPackage()
//	oldPackage := &mysql.PluginPackage{
//		AppUUID: h.req.OldAppUUID,
//		Version: h.req.OldAppVersion,
//	}
//	exist2, err := pkgModel.Exist(oldPackage)
//	if err != nil {
//		log.ErrorDetails(errors.Trace(err))
//		return errors.PluginUpgradeError(errors.ServerError)
//	}
//	if !exist2 {
//		return errors.PluginUpgradeError(errors.PackageNotFound)
//	}
//	h.oldPackage = oldPackage
//
//	// 判断oldInstance是否存在
//	instanceModel := mysql.ModelPluginInstance()
//	oldInstance := &mysql.PluginInstance{
//		AppUUID: h.req.OldAppUUID,
//		Version: h.req.OldAppVersion,
//	}
//	exist1, err := instanceModel.Exist(oldInstance)
//	if err != nil {
//		log.ErrorDetails(errors.Trace(err))
//		return errors.PluginUpgradeError(errors.ServerError)
//	}
//	if !exist1 {
//		return errors.PluginUpgradeError(errors.InstanceNotFound)
//	}
//	h.oldInstance = oldInstance
//
//	// 获取yaml文件
//	config, err := oldPackage.LoadYamlConfig()
//	if err != nil {
//		return errors.Trace(err)
//	}
//	h.oldConfig = config
//
//	// 对比host
//	host := handler.GetHost(h.oldInstance.InstanceUUID)
//	if host == nil {
//		return errors.PluginUpgradeError(errors.HostNotFound)
//	}
//	if host.GetInfo().Language != h.newConfig.Language {
//		return errors.PluginUpgradeError(errors.LanguageDisMatch)
//	}
//	if host.GetInfo().LanguageVersion != h.newConfig.LanguageVersion {
//		return errors.PluginUpgradeError(errors.LanguageVersionDisMatch)
//	}
//	if host.GetInfo().Version != h.newConfig.HostVersion {
//		return errors.PluginUpgradeError(errors.HostLanguageVersionDisMatch)
//	}
//	return nil
//}
//
//func (h *UpgradeHelper) upgrade() error {
//	oldVersion := &protocol.PluginDescriptor{
//		ApplicationID:      h.oldConfig.AppUUID,
//		Name:               h.oldConfig.Name,
//		Language:           h.oldConfig.Language,
//		LanguageVersion:    message.VersionString2Pb(h.oldConfig.LanguageVersion),
//		ApplicationVersion: message.VersionString2Pb(h.oldConfig.Version),
//		HostVersion:        message.VersionString2Pb(h.oldConfig.HostVersion),
//		MinSystemVersion:   message.VersionString2Pb(h.oldConfig.MinSystemVersion),
//	}
//	er := handler.UpgradePlugin(h.newConfig.AppUUID, h.oldInstance.InstanceUUID, h.newConfig.Name,
//		h.newConfig.Language, h.newConfig.LanguageVersion, h.newConfig.Version, oldVersion)
//	if er != nil {
//		log.PEDetails(er)
//		return errors.PluginDisableError(er.Error() + " " + er.Msg())
//	}
//	return nil
//}
