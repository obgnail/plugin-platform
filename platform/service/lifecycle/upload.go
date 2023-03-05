package lifecycle

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/obgnail/plugin-platform/common/errors"
	"github.com/obgnail/plugin-platform/common/log"
	"github.com/obgnail/plugin-platform/platform/model/mysql"
	"github.com/obgnail/plugin-platform/platform/pool/plugin_pool"
	"github.com/obgnail/plugin-platform/platform/service/types"
	"github.com/obgnail/plugin-platform/platform/service/utils"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"path/filepath"
)

type UploadResponse struct {
	mysql.PluginPackage
	IsUpgrade bool `json:"upgrade"`
}

func Upload(ctx *gin.Context) (ret gin.H, err error) {
	helper := &UploadHelper{ctx: ctx}
	if err = helper.checkFile(); err != nil {
		return ret, errors.Trace(err)
	}
	if err = helper.parseConfigYaml(); err != nil {
		return ret, errors.Trace(err)
	}
	if err = helper.checkInstanceRunning(); err != nil {
		return ret, errors.Trace(err)
	}
	upgrade, oldPackage, err := helper.checkUpgrade()
	if err != nil {
		return ret, errors.Trace(err)
	}
	if upgrade {
		if err = helper.upgrade(oldPackage); err != nil {
			return ret, errors.Trace(err)
		}
	} else {
		if err = helper.upload(); err != nil {
			return ret, errors.Trace(err)
		}
	}
	resp, err := helper.newPackage(upgrade)
	if err != nil {
		return ret, errors.Trace(err)
	}
	return gin.H{"data": resp}, nil
}

type UploadHelper struct {
	ctx        *gin.Context
	fileHeader *multipart.FileHeader
	file       multipart.File

	pluginConfig *plugin_pool.PluginConfig // config.yaml
	appUUID      string
	version      string
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

	var pluginConfig = new(plugin_pool.PluginConfig)
	if err = yaml.Unmarshal([]byte(configYaml), pluginConfig); err != nil {
		log.Warn("Unmarshal err: %+v", errors.Trace(err))
		return errors.PluginUploadError(errors.PluginConfigFileParseFailed)
	}

	h.pluginConfig = pluginConfig
	h.appUUID = h.pluginConfig.Service.AppUUID
	h.version = h.pluginConfig.Service.Version
	return nil
}

func (h *UploadHelper) checkInstanceRunning() error {
	instance := &mysql.PluginInstance{AppUUID: h.appUUID}
	instanceExist, err := mysql.ModelPluginInstance().Exist(instance)
	if err != nil {
		return errors.PluginUploadError(errors.ServerError)
	}
	if instanceExist {
		// 必须先停下实例才能升级
		return errors.PluginUploadError(errors.PluginAlreadyRunning)
	}
	return nil
}

func (h *UploadHelper) checkUpgrade() (upgrade bool, oldPackage *mysql.PluginPackage, err error) {
	oldPackage = &mysql.PluginPackage{AppUUID: h.appUUID}
	exist, err := mysql.ModelPluginPackage().Exist(oldPackage)
	if err != nil {
		log.ErrorDetails(errors.Trace(err))
		err = errors.PluginUploadError(errors.ServerError)
		return
	}
	if exist {
		compare, err := utils.PluginVersionCompare(h.version, oldPackage.Version)
		if err != nil {
			log.ErrorDetails(errors.Trace(err))
			return false, nil, errors.PluginUploadError(errors.ServerError)
		}
		// 上传的包版本大于当前安装包的版本，符合升级条件
		return compare == types.VersionMore, oldPackage, nil
	}

	return
}

func (h *UploadHelper) upgrade(oldPackage *mysql.PluginPackage) error {
	if err := oldPackage.RealDelete(oldPackage.Id, oldPackage); err != nil {
		log.ErrorDetails(errors.Trace(err))
		return errors.PluginUploadError(errors.ServerError)
	}
	return nil
}

func (h *UploadHelper) upload() error {
	dirPath := utils.GetPluginDir(h.appUUID, h.version)
	if err := utils.SaveDecompressedFiles(h.fileHeader, dirPath); err != nil {
		log.ErrorDetails(errors.Trace(err))
		return errors.PluginUploadError(errors.FileUnzipFailed)
	}
	return nil
}

func (h *UploadHelper) newPackage(upgrade bool) (*UploadResponse, error) {
	m := mysql.ModelPluginPackage()
	one := &mysql.PluginPackage{
		AppUUID: h.appUUID,
		Name:    h.pluginConfig.Name,
		Size:    h.fileHeader.Size,
		Version: h.version,
	}
	if err := m.New(one); err != nil {
		log.ErrorDetails(errors.Trace(err))
		return nil, errors.PluginUploadError(errors.ServerError)
	}

	resp := &UploadResponse{
		PluginPackage: *one,
		IsUpgrade:     upgrade,
	}
	return resp, nil
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
