package protocol

import (
	"strconv"
	"strings"
)

func SplitVersion(version string) *Version {
	versionSlice := strings.Split(version, ".")
	major, _ := strconv.ParseInt(versionSlice[0], 10, 32)
	minor, _ := strconv.ParseInt(versionSlice[1], 10, 32)
	patch, _ := strconv.ParseInt(versionSlice[2], 10, 32)
	return &Version{Major: int32(major), Minor: int32(minor), Revision: int32(patch)}
}
