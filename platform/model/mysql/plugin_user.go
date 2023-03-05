package mysql

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

func ModelPluginUser() *PluginUser {
	var m = new(PluginUser)
	m.Child = m
	return m
}

type PluginUser struct {
	BaseModel
	UserUUID     string `gorm:"user_uuid" json:"user_uuid"`
	AppUUID      string `gorm:"app_uuid" json:"app_uuid"`
	InstanceUUID string `gorm:"instance_uuid" json:"instance_uuid"`
	Name         string `gorm:"name" json:"name"`
	Email        string `gorm:"email" json:"email"`
}

func (u *PluginUser) tableName() string {
	return "plugin_user"
}

func (u *PluginUser) Uninstall(db *gorm.DB, appId string, instanceId string) error {
	err := db.Table(u.tableName()).
		Where("app_uuid = ? and instance_uuid = ? ", appId, instanceId).
		Updates(map[string]interface{}{
			"deleted": true,
		}).Error
	if err != nil {
		return err
	}
	return nil
}

func NewUserUUID(appUUID, instanceUUID string) string {
	return fmt.Sprintf("plugin_%s_%s", appUUID, instanceUUID)
}

func NewUserEmail(appUUID, instanceUUID string) string {
	return fmt.Sprintf("plugin_%s_%s@email.com", appUUID, instanceUUID)
}
