package lifecycle

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/obgnail/plugin-platform/common/errors"
	"github.com/obgnail/plugin-platform/common/log"
	"github.com/obgnail/plugin-platform/platform/model/mysql"
	"github.com/obgnail/plugin-platform/platform/service/common"
	"github.com/obgnail/plugin-platform/platform/service/types"
	"github.com/obgnail/plugin-platform/platform/service/utils"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"path/filepath"
)

type UploadResponse struct {
	AppUUID      string `json:"app_uuid"`
	Version      string `json:"version"`
	InstanceUUID string `json:"instance_uuid"`
	Name         string `json:"name"`
	NewVersion   string `json:"new_version"`
	IsUpgrade    bool   `json:"is_upgrade"`
}

func Upload(ctx *gin.Context) (ret gin.H, err error) {
	helper := &UploadHelper{ctx: ctx}
	if err = helper.checkFile(); err != nil {
		return ret, errors.Trace(err)
	}
	if err = helper.parseConfigYaml(); err != nil {
		return ret, errors.Trace(err)
	}
	instance, exist, err := helper.getInstance()
	if err != nil {
		return ret, errors.Trace(err)
	}
	// instance不存在
	if !exist {
		return helper.uploadNonexistent()
	}
	// 覆盖安装(旧插件必须没有install过)
	if instance.Status == common.PluginStatusUploaded {
		return helper.overlayExistent(instance)
	} else {
		// 升级
		return helper.upgrade(instance)
	}
}

type UploadHelper struct {
	ctx        *gin.Context
	fileHeader *multipart.FileHeader
	file       multipart.File

	cfg *common.PluginConfig // config.yaml
}

func (h *UploadHelper) uploadNonexistent() (ret gin.H, err error) {
	if err = h.unpack(); err != nil {
		return ret, errors.Trace(err)
	}
	instance, err := h.newInstance()
	if err != nil {
		return ret, errors.Trace(err)
	}
	resp := &UploadResponse{
		AppUUID:      instance.AppUUID,
		Version:      instance.Version,
		NewVersion:   "",
		IsUpgrade:    false,
		InstanceUUID: instance.InstanceUUID,
		Name:         instance.Name,
	}
	return gin.H{"data": resp}, nil
}

func (h *UploadHelper) overlayExistent(instance *mysql.PluginInstance) (ret gin.H, err error) {
	if err = mysql.ModelPluginInstance().Delete(instance.Id, instance); err != nil {
		log.ErrorDetails(errors.Trace(err))
		return ret, errors.PluginDisableError(errors.ServerError)
	}
	if ret, err = h.uploadNonexistent(); err != nil {
		log.ErrorDetails(errors.Trace(err))
		return ret, errors.Trace(err)
	}
	return ret, nil
}

func (h *UploadHelper) upgrade(instance *mysql.PluginInstance) (ret gin.H, err error) {
	upgrade, err := h.checkUpgrade(instance)
	if err != nil {
		return ret, errors.Trace(err)
	}
	if upgrade {
		if err = h.upgradeFile(); err != nil {
			return ret, errors.Trace(err)
		}
	}
	//resp, err := h.newInstance(upgrade)
	//if err != nil {
	//	return ret, errors.Trace(err)
	//}
	// TODO
	return gin.H{"data": "00"}, nil
}

func (h *UploadHelper) checkFile() error {
	err := CheckMaxFileSize(h.ctx, types.PluginFileMaxSize)
	if err != nil {
		return errors.PluginUploadError(errors.FileTooLarge)
	}
	h.file, h.fileHeader, err = h.ctx.Request.FormFile("file")
	if err != nil {
		return errors.PluginUploadError(errors.FileUploadFailed)
	}
	if filepath.Ext(h.fileHeader.Filename) != types.PluginFileExt {
		return errors.PluginUploadError(errors.InvalidFileExt)
	}
	return nil
}

func (h *UploadHelper) parseConfigYaml() error {
	configYaml, err := utils.GetYamlFromFile(h.fileHeader, h.file)
	if err != nil {
		return errors.PluginUploadError(errors.PluginConfigFileParseFailed)
	}

	var pluginConfig = new(common.PluginConfig)
	if err = yaml.Unmarshal([]byte(configYaml), pluginConfig); err != nil {
		log.Warn("Unmarshal err: %+v", errors.Trace(err))
		return errors.PluginUploadError(errors.PluginConfigFileParseFailed)
	}

	h.cfg = pluginConfig
	return nil
}

func (h *UploadHelper) getInstance() (instance *mysql.PluginInstance, exist bool, err error) {
	instance = &mysql.PluginInstance{AppUUID: h.cfg.AppUUID}
	exist, err = mysql.ModelPluginInstance().Exist(instance)
	if err != nil {
		log.ErrorDetails(errors.Trace(err))
		err = errors.PluginUploadError(errors.ServerError)
	}
	return
}

func (h *UploadHelper) checkUpgrade(instance *mysql.PluginInstance) (upgrade bool, err error) {
	compare, err := utils.PluginVersionCompare(h.cfg.Version, instance.Version)
	if err != nil {
		log.ErrorDetails(errors.Trace(err))
		return false, errors.PluginUploadError(errors.ServerError)
	}
	// 上传的包版本大于当前安装包的版本，符合升级条件
	return compare == types.VersionMore, nil
}

func (h *UploadHelper) unpack() error {
	dirPath := utils.GetPluginDir(h.cfg.AppUUID, h.cfg.Version)
	if err := utils.SaveDecompressedFiles(h.fileHeader, dirPath); err != nil {
		log.ErrorDetails(errors.Trace(err))
		return errors.PluginUploadError(errors.FileUnzipFailed)
	}
	return nil
}

func (h *UploadHelper) upgradeFile() error {
	if err := h.unpack(); err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (h *UploadHelper) newInstance() (*mysql.PluginInstance, error) {
	apis, err := json.Marshal(h.cfg.Apis)
	if err != nil {
		log.ErrorDetails(errors.Trace(err))
		return nil, errors.PluginUploadError(errors.LoadYamlConfigFailed)
	}

	one := &mysql.PluginInstance{
		AppUUID:      h.cfg.AppUUID,
		InstanceUUID: utils.NewInstanceUUID(),
		Name:         h.cfg.Name,
		Version:      h.cfg.Version,
		Description:  h.cfg.Description,
		Contact:      h.cfg.Contact,
		Status:       common.PluginStatusUploaded,
		Apis:         string(apis),
	}
	if err := mysql.ModelPluginInstance().New(one); err != nil {
		log.ErrorDetails(errors.Trace(err))
		return nil, errors.PluginUploadError(errors.ServerError)
	}
	return one, nil
}

func CheckMaxFileSize(context *gin.Context, maxSize int64) error {
	var c = context

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)

	buff, err := c.GetRawData()
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(buff)
	c.Request.Body = ioutil.NopCloser(buf)
	return nil
}
