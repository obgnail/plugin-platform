package types

const (
	VersionLess  = 0
	VersionEqual = 1
	VersionMore  = 2

	VersionLen = 3
)

const (
	PluginFileMaxSize = 512 * 1024 * 1024

	PluginFileExt          = ".zip"
	PluginFileWebDir       = "web"
	PluginFileBackDir      = "backend"
	PluginServerDir        = "server"
	PluginFileWebDistDir   = "dist"
	PluginFileWebBuildDir  = "build"
	PluginFileConfigDir    = "config"
	PluginFileWorkspaceDir = "workspace"
	PluginUpgradeDir       = "upgrade"
	PluginFileConfigYaml   = "plugin.yaml"
	PluginSoFile           = "backend.so"
	PluginFileUpgradeYaml  = "upgrade.yaml"
	PluginFileLogo         = "logo.svg"
)
