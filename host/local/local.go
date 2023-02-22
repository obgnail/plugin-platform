package local

import (
	"github.com/obgnail/plugin-platform/host/config"
	"github.com/obgnail/plugin-platform/host/handler"
)

func StartHost() {
	id := config.StringOrPanic("runtime_id")
	name := config.StringOrPanic("runtime_name")
	addr := config.StringOrPanic("platform_address")
	lang := config.StringOrPanic("runtime_language")
	hostVersion := config.StringOrPanic("runtime_version")
	minSysVersion := config.StringOrPanic("runtime_min_system_version")
	langVersion := config.StringOrPanic("runtime_language_version")
	isLocal := true

	handler.New(id, name, addr, lang, hostVersion, minSysVersion, langVersion, isLocal)
}
