package mysql

import (
	"github.com/jinzhu/gorm"
	"strings"
)

func ModelPluginConfig() *PluginConfig {
	var m = new(PluginConfig)
	m.Child = m
	return m
}

type PluginConfig struct {
	BaseModel
	AppUUID      string `gorm:"COLUMN:app_uuid" json:"app_uuid"`
	InstanceUUID string `gorm:"COLUMN:instance_uuid" json:"instance_uuid"`
	Label        string `gorm:"COLUMN:label" json:"label"`
	Key          string `gorm:"COLUMN:arg_key" json:"arg_key"`
	Value        string `gorm:"COLUMN:arg_value" json:"arg_value"`
	Type         int64  `gorm:"COLUMN:type" json:"type"` // COMBO: 1,INPUT: 2, TEXT: 3, SELECT: 4,CHECKBOX: 5,BUTTON: 6,
	Required     bool   `gorm:"COLUMN:required" json:"required"`
}

func (c *PluginConfig) tableName() string {
	return "plugin_config"
}

func ConvertConfigType(Type string) int64 {
	switch strings.ToUpper(Type) {
	case "COMBO":
		return 1
	case "INPUT":
		return 2
	case "TEXT":
		return 3
	case "SELECT":
		return 4
	case "CHECKBOX":
		return 5
	case "BUTTON":
		return 6
	default:
		return 0
	}
}

func (c *PluginConfig) Uninstall(db *gorm.DB, appId string, instanceId string) error {
	err := db.Table(c.tableName()).
		Where("app_uuid = ? and instance_uuid = ?", appId, instanceId).
		Updates(map[string]interface{}{
			"deleted": true,
		}).Error
	if err != nil {
		return err
	}
	return nil
}
