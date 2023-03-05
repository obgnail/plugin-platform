package errors

const (

	// 文件过大
	FileTooLarge = "FileTooLarge"

	// 文件上传失败
	FileUploadFailed = "FileUploadFailed"

	// 文件格式不正确
	FileMalformed = "FileMalformed"

	// 文件解压失败
	FileUnzipFailed = "FileUnzipFailed"

	// 加载配置文件失败
	LoadYamlConfigFailed = "LoadYamlConfig"

	// 文件元数据错误
	ErrorMetaData = "ErrorMetaData"

	// 保存到数据库失败
	SaveToDBFailed = "SaveToDBFailed"

	// 插件文件不存在
	FileNoExist = "FileNoExist"

	// 插件包已存在
	FileAlreadyExist = "FileAlreadyExist"

	// 未找到该实例
	InstanceNotFound = "InstanceNotFound"

	// 未找到package
	PackageNotFound = "PackageNotFound"

	// 已经卸载
	PluginAlreadyUninstall = "PluginAlreadyUninstall"

	// 插件已安装
	PluginAlreadyInstall = "PluginAlreadyInstall"

	// 插件配置文件解析失败
	PluginConfigFileParseFailed = "PluginConfigFileParseFailed"
)

const (
	InstanceUUID = "instance_uuid"
	AppUUID      = "app_uuid"
	AppVersion   = "app_version"
)
