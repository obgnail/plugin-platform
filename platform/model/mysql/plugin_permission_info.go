package mysql

func ModelPluginPermissionInfo() *PluginPermissionInfo {
	var m = new(PluginPermissionInfo)
	m.Child = m
	return m
}

type PluginPermissionInfo struct {
	BaseModel
	InstanceUUID    string `gorm:"COLUMN:instance_uuid" json:"instance_uuid"`
	PermissionName  string `gorm:"COLUMN:permission_name" json:"permission_name"`
	PermissionField string `gorm:"COLUMN:permission_field" json:"permission_field"`
	PermissionDesc  string `gorm:"COLUMN:permission_desc" json:"permission_desc"`
	PermissionID    int    `gorm:"COLUMN:permission_id" json:"permission_id"`
}

func (i *PluginPermissionInfo) tableName() string {
	return "plugin_permission_info"
}
