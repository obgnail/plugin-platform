package mysql

import (
	"github.com/obgnail/plugin-platform/common/errors"
	"github.com/obgnail/plugin-platform/platform/pool/plugin_pool"
	"github.com/obgnail/plugin-platform/platform/service/utils"
	"gopkg.in/yaml.v2"
)

func ModelPluginPackage() *PluginPackage {
	var m = new(PluginPackage)
	m.Child = m
	return m
}

type PluginPackage struct {
	BaseModel
	AppUUID string `gorm:"app_uuid" json:"app_uuid"`
	Name    string `gorm:"name" json:"name"`
	Size    int64  `gorm:"size" json:"size"`
	Version string `gorm:"version" json:"version"`
}

func (p *PluginPackage) tableName() string {
	return "plugin_package"
}

func (p *PluginPackage) LoadYamlConfig() (*plugin_pool.PluginConfig, error) {
	yamlPath := utils.GetPluginConfigPath(p.AppUUID, p.Version)
	res, err := utils.ReadFile(yamlPath)
	if err != nil {
		return nil, errors.Trace(err)
	}

	var pluginConfig = new(plugin_pool.PluginConfig)
	if err := yaml.Unmarshal(res, pluginConfig); err != nil {
		return nil, errors.Trace(err)
	}
	return pluginConfig, nil
}
