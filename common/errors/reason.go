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

	// 插件文件不存在
	FileNoExist = "FileNoExist"

	// 插件包已存在
	FileAlreadyExist = "FileAlreadyExist"

	// 插件已安装
	PluginAlreadyInstall = "PluginAlreadyInstall"

	// 插件配置文件解析失败
	PluginConfigFileParseFailed = "PluginConfigFileParseFailed"
)

const (
	InstanceUUID = "instance_uuid"
	AppUUID      = "app_uuid"
)
