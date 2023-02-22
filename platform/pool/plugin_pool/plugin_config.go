package plugin_pool

const (
	Addition = "addition"
	Replace  = "replace"
	Prefix   = "prefix"
	Suffix   = "suffix"
	External = "external"
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
	Status           int           `yaml:"status" json:"status"`
	Config           []*Config     `yaml:"config" json:"config"`         //插件自己的配置
	Permission       []*Permission `yaml:"permission" json:"permission"` // 插件自身的权限
	Path             string        `yaml:"path" json:"path"`             // 非config使用，nginx加载前端资源的路径
}

// 插件自定义的路由
type Api struct {
	Type     string   `yaml:"type" json:"type"` // add pre suf replace external(TODO)
	Methods  []string `yaml:"methods" json:"methods"`
	Url      string   `yaml:"url" json:"url"`
	Function string   `yaml:"function" json:"function"` // handler函数
}

type Ability struct {
	Id          string                 `yaml:"id" json:"id"`
	Name        string                 `yaml:"name" json:"name"`
	AbilityType string                 `yaml:"abilityType" json:"abilityType"` // bang-api支持的能力（）
	Version     string                 `yaml:"version" json:"version"`         // bang-api支持能力的版本
	Label       string                 `yaml:"label" json:"label"`
	Function    map[string]interface{} `yaml:"function" json:"function"` // handler函数组
	Setting     map[string]interface{} `yaml:"setting" json:"setting"`   // 可以当成没有
	Config      []interface{}          `yaml:"config" json:"config"`     // 每个能力都有自己的配置，
}

type PluginConfig struct {
	*Service      `yaml:"service" json:"service"`
	Apis          []*Api      `yaml:"apis" json:"apis"`
	Modules       interface{} `yaml:"modules" json:"modules"`     // 前端使用，用于生成前端插槽
	Abilities     []*Ability  `yaml:"abilities" json:"abilities"` // 插件能力的描述
	WebServiceUrl string      `json:"web_service_url,omitempty"`  // 本地开发时，访问插件前端资源的路由
}