package common

type RouterType = string

// addition和external的异同:
//   前者是加在主系统里的,补充主系统的接口。
//   后者是插件提供一个暴露到到程序外面的接口,提供对外的http服务。
//   二者是根据业务来区分的,最终url不同,调用的函数也不同:前者是自定义,后者强制OnExternalHttpRequest
//   当然,二者底层实现没有区别,所以硬要混着用也行:)
const (
	RouterTypeAddition RouterType = "addition"
	RouterTypeExternal RouterType = "external"
	RouterTypeReplace  RouterType = "replace"
	RouterTypePrefix   RouterType = "prefix"
	RouterTypeSuffix   RouterType = "suffix"
)

type StatusType = int

const (
	PluginStatusUploaded    StatusType = 1 // 已上传
	PluginStatusRunning     StatusType = 2 // 已启用
	PluginStatusStopping    StatusType = 3 // 已安装/已停用
	PluginStatusUnavailable StatusType = 4 // 不可用,插件升级后,旧插件会被设置为不可用
)

type Config struct {
	Key      string `yaml:"key" json:"key"`
	Value    string `yaml:"value" json:"value"`
	Type     string `yaml:"type" json:"type"`
	Required bool   `yaml:"required" json:"required"`
	Label    string `yaml:"label" json:"label"`
}

type Permission struct {
	Name  string `yaml:"name" json:"name"`
	Field string `yaml:"field" json:"field"`
	Desc  string `yaml:"desc" json:"desc"`
}

type Service struct {
	AppUUID          string        `yaml:"app_id" json:"app_id"`           // 生成插件zip文件时自动生成appUUID
	InstanceUUID     string        `yaml:"instance_id" json:"instance_id"` // 内部流转的时候使用的，技术上是可以代替appUUID。引入的原因是appUUID不可靠
	Name             string        `yaml:"name" json:"name"`
	Version          string        `yaml:"version" json:"version"`
	Description      string        `yaml:"description" json:"description"`
	Language         string        `yaml:"language" json:"language"` // 开发语言
	LanguageVersion  string        `yaml:"language_version" json:"language_version"`
	Icon             string        `json:"icon"`
	HostVersion      string        `yaml:"host_version" json:"host_version"`             // 生成zip就已经有Host
	MinSystemVersion string        `yaml:"min_system_version" json:"min_system_version"` // 能允许插件的主程序的最小版本
	Contact          string        `yaml:"contact" json:"contact"`                       // 开发者的联系方式
	Status           StatusType    `yaml:"status" json:"status"`
	Config           []*Config     `yaml:"config" json:"config"`         //插件自己的配置
	Permission       []*Permission `yaml:"permission" json:"permission"` // 插件自身的权限
	Path             string        `yaml:"path" json:"path"`             // 非config使用,nginx加载前端资源的路径
}

// Api 插件自定义的路由
type Api struct {
	Type     RouterType `yaml:"type" json:"type"` // addition、replace、suffix、prefix、external
	Methods  []string   `yaml:"methods" json:"methods"`
	Url      string     `yaml:"url" json:"url"`
	Function string     `yaml:"function" json:"function"` // handler函数
}

type Ability struct {
	Id          string            `yaml:"id" json:"id"`
	Name        string            `yaml:"name" json:"name"`
	AbilityType string            `yaml:"abilityType" json:"abilityType"` // 主系统支持的能力
	Version     string            `yaml:"version" json:"version"`         // 主系统支持能力的版本
	Label       string            `yaml:"label" json:"label"`
	Function    map[string]string `yaml:"function" json:"function"` // handler函数组
	Setting     map[string]string `yaml:"setting" json:"setting"`   // 可以当成没有
	Config      []interface{}     `yaml:"config" json:"config"`     // 每个能力都有自己的配置，
}

type PluginConfig struct {
	*Service      `yaml:"service" json:"service"`
	Apis          []*Api      `yaml:"apis" json:"apis"`
	Modules       interface{} `yaml:"modules" json:"modules"`     // 前端使用，用于生成前端插槽
	Abilities     []*Ability  `yaml:"abilities" json:"abilities"` // 插件能力的描述
	WebServiceUrl string      `json:"web_service_url,omitempty"`  // 本地开发时，访问插件前端资源的路由
}
