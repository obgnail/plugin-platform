package mysql

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/obgnail/plugin-platform/common/errors"
	"github.com/obgnail/plugin-platform/platform/service/types"
	"github.com/obgnail/plugin-platform/platform/service/utils"
	"gopkg.in/yaml.v2"
	"reflect"
)

func ModelPluginInstance() *PluginInstance {
	var m = new(PluginInstance)
	m.Child = m
	return m
}

type PluginInstance struct {
	BaseModel
	AppUUID      string `gorm:"app_uuid" json:"app_uuid"`
	InstanceUUID string `gorm:"instance_uuid" json:"instance_uuid"`
	Name         string `gorm:"name" json:"name"`
	Version      string `gorm:"version" json:"version"`
	Description  string `gorm:"description" json:"description"`
	Contact      string `gorm:"contact" json:"contact"`
	Status       int    `gorm:"status" json:"status"`
	Apis         string `gorm:"apis" json:"apis"`
}

func (i *PluginInstance) tableName() string {
	return "plugin_instance"
}

func (i *PluginInstance) GetConfig() *types.PluginConfig {
	s := &types.Service{
		AppUUID:      i.AppUUID,
		InstanceUUID: i.InstanceUUID,
		Name:         i.Name,
		Version:      i.Version,
		Description:  i.Description,
		Status:       i.Status,
	}
	apis := make([]*types.Api, 0)
	_ = json.Unmarshal([]byte(i.Apis), &apis)
	cfg := &types.PluginConfig{
		Service: s,
		Apis:    apis,
	}
	return cfg
}

func (i *PluginInstance) LoadYamlConfig() (*types.PluginConfig, error) {
	yamlPath := utils.GetPluginConfigPath(i.AppUUID, i.Version)
	res, err := utils.ReadFile(yamlPath)
	if err != nil {
		return nil, errors.Trace(err)
	}

	var pluginConfig = new(types.PluginConfig)
	if err := yaml.Unmarshal(res, pluginConfig); err != nil {
		return nil, errors.Trace(err)
	}
	return pluginConfig, nil
}

func (i *PluginInstance) FirstInstance(appId string) (*PluginInstance, error) {
	var all = make([]*PluginInstance, 0)
	err := DB.Table(i.tableName()).
		Where("app_id = ? and deleted = ?", appId, false).
		Find(&all).Error
	if err != nil {
		return &PluginInstance{}, err
	}

	var min = new(PluginInstance)
	if len(all) > 0 {
		min = all[0]
	}

	for _, a := range all {
		if a.CreateTime < min.CreateTime {
			min = a
		}
	}
	return min, nil
}

type AppInstanceCount struct {
	AppId string
	Num   int64
}

func (b *AppInstanceCount) String() string {
	className := reflect.TypeOf(b).Elem().Name()
	s := "< " + className + "\n"

	e := reflect.ValueOf(b).Elem()
	for i := 0; i < e.NumField(); i++ {
		attr := e.Type().Field(i).Name
		sType := e.Type().Field(i).Type
		v := e.Field(i).Interface()
		s += fmt.Sprintf("%v: %v(%v)\n", attr, sType, v)
	}
	s += ">\n"
	return s
}

func (i *PluginInstance) GroupByAppId() ([]*AppInstanceCount, error) {
	var result = make([]*AppInstanceCount, 0)
	err := DB.Table(i.tableName()).
		Select("app_id as app_id, count(*) as num").
		Group("app_id").Where("deleted = ?", false).
		Find(&result).Error
	if err != nil {
		return result, err
	}

	return result, nil
}

// 获取除启用状态外的所有插件实例ID
func (i *PluginInstance) GetExcludeStartInstanceList() ([]*PluginInstance, error) {
	var result = make([]*PluginInstance, 0)
	err := DB.Table(i.tableName()).
		Select("*").
		Where("status != ?", 1).
		Find(&result).Error
	if err != nil {
		return result, err
	}

	return result, nil
}

func (i *PluginInstance) GetOnePluginInstance(orgUUID, teamUUID, appUUID, instanceUUID string) ([]*PluginInstance, error) {
	var result = make([]*PluginInstance, 0)
	err := DB.Table(i.tableName()).
		Select("*").
		Where("status = ?", 2).
		Where("deleted = ?", 1).
		Where("org_uuid = ?", orgUUID).
		Where("team_uuid = ?", teamUUID).
		Where("app_uuid = ?", appUUID).
		Where("instance_uuid = ?", instanceUUID).
		Find(&result).Error
	if err != nil {
		return result, err
	}

	return result, nil
}

func (i *PluginInstance) RealDeleteWithDBArg(db *gorm.DB, arg ModelInter) error {
	err := db.Table(i.Child.tableName()).
		Where(arg).
		Delete(i.Child).
		Error

	if err != nil {
		e := reflect.ValueOf(arg).Elem()
		className := e.Type().Name()
		err = fmt.Errorf("%s.Delete: %v", className, err)
		return err
	}
	return nil
}
