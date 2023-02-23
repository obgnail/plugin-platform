package protocol

import (
	"github.com/obgnail/plugin-platform/common/common_type"
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

func BuildErrorMessage(err common_type.PluginError) *ErrorMessage {
	if err != nil {
		return &ErrorMessage{
			Code:  int64(err.Code()),
			Error: err.Error(),
			Msg:   err.Msg(),
		}
	}

	return &ErrorMessage{}
}

func BuildResourceFileMessage(distinctMessage *PlatformMessage, resp *WorkspaceMessage_IOResponseMessage) {
	// workspaceMessage
	workspaceMessage := &WorkspaceMessage{}
	workspaceMessage.IOResponse = resp
	distinctMessage.Resource.Workspace = workspaceMessage
}
