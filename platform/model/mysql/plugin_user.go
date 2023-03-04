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
	UserUUID     string `gorm:"COLUMN:user_uuid" json:"user_uuid"`
	AppUUID      string `gorm:"COLUMN:app_uuid" json:"app_uuid"`
	InstanceUUID string `gorm:"COLUMN:instance_uuid" json:"instance_uuid"`
	Name         string `gorm:"COLUMN:name" json:"name"`
	Email        string `gorm:"COLUMN:email" json:"email"`
}

func (u *PluginUser) tableName() string {
	return "plugin_user"
}

func (u *PluginUser) Uninstall(db *gorm.DB, appId string, instanceId string, orgId string, teamId string) error {
	err := db.Table(u.tableName()).
		Where("app_uuid = ? and instance_uuid = ? and org_uuid = ? and team_uuid = ?", appId, instanceId, orgId, teamId).
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
