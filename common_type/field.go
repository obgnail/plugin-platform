package common

type Field interface {
	ItemsAddProject(*ItemsAddProjectReq) *ItemsAddResp
	ItemsAddProduct(*ItemsAddProductReq) *ItemsAddResp
	FieldsAdd(*FieldsAddReq) *FieldsAddResp
	AddGroupField(*AddGroupFieldReq) *AddGroupFieldResp
	UpdateFieldOption(*UpdateFieldOptionReq) PluginError
}

type ItemsAddProjectReq struct {
	FieldType 		string
	Name      		string
	ItemType  		string
	Pool      		string
	ContextType  	string
}

type ItemsAddProductReq struct {
	FieldType 		string
	Name      		string
	ItemType  		string
	Pool      		string
	ContextType  	string
	Require         bool
}

type ItemsAddResp struct {
	UUID  string
	Error PluginError
}

type FieldsAddReq struct {
	Name         string
	Type         int64
	Renderer     int64
	FilterOption int64
	SearchOption int64
}

type FieldsAddResp struct {
	UUID  string
	Error PluginError
}

type AddGroupFieldReq struct {
	FieldGroups FieldGroup
}

type FieldGroup struct {
	ObjectType string
	Name       string
	Relations  []Relation
}

type Relation struct {
	FieldUUID       string
	FieldParentUUID string
	Position        int64
}

type AddGroupFieldResp struct {
	Error PluginError
	UUIDs string
}

type UpdateFieldOptionReq struct {
	Options []ScriptFieldOption
}

type ScriptFieldOption struct {
	TeamUUID   string
	FieldUUID  string
	UUID       string
	Value      string
	ObjectType int64
}
