package message_utils

import (
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/config"
	"github.com/obgnail/plugin-platform/common/protocol"
	"strconv"
	"strings"
)

func VersionString2Pb(version string) *protocol.Version {
	versionSlice := strings.Split(version, ".")
	major, _ := strconv.ParseInt(versionSlice[0], 10, 32)
	minor, _ := strconv.ParseInt(versionSlice[1], 10, 32)
	patch, _ := strconv.ParseInt(versionSlice[2], 10, 32)
	return &protocol.Version{Major: int32(major), Minor: int32(minor), Revision: int32(patch)}
}

func NewProtocolVersion(version common_type.IVersion) *protocol.Version {
	return &protocol.Version{Major: int32(version.Major()), Minor: int32(version.Minor()), Revision: int32(version.Revision())}
}

func VersionPb2String(version *protocol.Version) string {
	major := version.GetMajor()
	minor := version.GetMinor()
	revision := version.GetRevision()

	var Major, Minor, Revision string
	if major != 0 {
		Major = strconv.FormatInt(int64(major), 10)
	} else {
		Major = "0"
	}

	if minor != 0 {
		Minor = strconv.FormatInt(int64(minor), 10)
	} else {
		Minor = "0"
	}

	if revision != 0 {
		Revision = strconv.FormatInt(int64(revision), 10)
	} else {
		Revision = "0"
	}

	applicationVersion := strings.Join([]string{Major, Minor, Revision}, ".")
	return applicationVersion
}

func IsLocalDB(dbName string) bool {
	return dbName == config.String("platform.mysql_user_plugin_db_name", "plugins")
}
