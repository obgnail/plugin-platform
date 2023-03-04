package mysql

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/obgnail/plugin-platform/platform/pool/plugin_pool"
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
	CreateTime   string `gorm:"create_time" json:"create_time"`
	UpdateTime   string `gorm:"update_time" json:"update_time"`
	Deleted      int    `gorm:"deleted" json:"deleted"`
}

func (i *PluginInstance) tableName() string {
	return "plugin_instance"
}

func (i *PluginInstance) GetConfig() *plugin_pool.PluginConfig {
	s := &plugin_pool.Service{
		AppUUID:      i.AppUUID,
		InstanceUUID: i.InstanceUUID,
		Name:         i.Name,
		Version:      i.Version,
		Description:  i.Description,
		Status:       i.Status,
	}
	apis := make([]*plugin_pool.Api, 0)
	_ = json.Unmarshal([]byte(i.Apis), &apis)
	cfg := &plugin_pool.PluginConfig{
		Service: s,
		Apis:    apis,
	}
	return cfg
}

//func (i *PluginInstance) Uninstall(appId string, instanceId string, orgId string, teamId string) error {
//	callback := func(tx *gorm.DB) error {
//		err := tx.Table(i.tableName()).
//			Where("app_uuid = ? and instance_uuid = ? and org_uuid = ? and team_uuid = ?", appId, instanceId, orgId, teamId).
//			Updates(map[string]interface{}{
//				"deleted": true,
//			}).Error
//		if err != nil {
//			return err
//		}
//
//		configModel := ModelPluginConfig()
//		if err = configModel.Uninstall(tx, appId, instanceId, orgId, teamId); err != nil {
//			return err
//		}
//
//		authModel := ModelPluginAuthDesc()
//		if err = authModel.Uninstall(tx, appId, instanceId, orgId, teamId); err != nil {
//			return err
//		}
//
//		permissionInfoModel := ModelPluginPermissionInfo()
//		if err = permissionInfoModel.Uninstall(tx, instanceId, orgId, teamId); err != nil {
//			return err
//		}
//
//		one := &PluginUser{
//			InstanceUUID: instanceId,
//			TeamUUID:     teamId,
//		}
//		m := ModelPluginUser()
//		if err = m.One(one); err != nil && err != RecordNotFound {
//			return err
//		}
//		if err = m.RealDeleteWithDB(tx, one.Id, one); err != nil {
//			return err
//		}
//
//		propertyModel := ModelPluginDataProperty()
//		if err = propertyModel.DelAbilityProperties(tx, appId, instanceId, orgId, teamId); err != nil {
//			return errors.Trace(err)
//		}
//
//		return nil
//	}
//	return DB.Transaction(callback)
//}

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
