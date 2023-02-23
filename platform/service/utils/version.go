package utils

import (
	"github.com/obgnail/plugin-platform/common/errors"
	"github.com/obgnail/plugin-platform/platform/service/types"
	"strconv"
	"strings"
)

var MalformedVersionError = errors.New("Malformed version string")

func PluginVersionCompare(version1, version2 string) (result int, err error) {
	arr1 := strings.Split(version1, ".")
	if len(arr1) != types.VersionLen {
		return 0, MalformedVersionError
	}
	arr2 := strings.Split(version2, ".")
	if len(arr2) != types.VersionLen {
		return 0, MalformedVersionError
	}

	for i := 0; i < types.VersionLen; i++ {
		num1, err := strconv.Atoi(arr1[i])
		if err != nil {
			return 0, errors.Trace(err)
		}
		num2, err := strconv.Atoi(arr2[i])
		if err != nil {
			return 0, errors.Trace(err)
		}

		if num1 < num2 {
			return types.VersionLess, nil
		} else if num1 > num2 {
			return types.VersionMore, nil
		}
	}
	return types.VersionEqual, nil
}
