package work_space

import (
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/message_utils"
	"github.com/obgnail/plugin-platform/common/protocol"
)

type WorkSpace struct {
	source     *protocol.PlatformMessage
	distinct   *protocol.PlatformMessage
	opType     protocol.WorkspaceMessage_IOOperationType
	AppID      string
	instanceID string
}

func NewWorkSpace(source, distinct *protocol.PlatformMessage) *WorkSpace {
	w := &WorkSpace{
		source:     source,
		distinct:   distinct,
		opType:     source.GetResource().GetWorkspace().GetIORequest().GetOperation(),
		AppID:      source.GetResource().GetSender().GetApplication().GetApplicationID(),
		instanceID: source.GetResource().GetSender().GetInstanceID(),
	}
	return w
}

func (w *WorkSpace) Execute() {
	var f = &SpaceOperation{AppID: w.AppID, InstanceID: w.instanceID}
	var err common_type.PluginError
	var ok bool
	var fileTree []string
	var fileByte []byte
	IOReqMsg := w.source.GetResource().GetWorkspace().GetIORequest()

	switch w.opType {
	case protocol.WorkspaceMessage_Create:
		err = f.CreateFile(IOReqMsg.GetFileName())
	case protocol.WorkspaceMessage_Rename:
		err = f.Rename(IOReqMsg.GetFileName(), IOReqMsg.GetNewFileName())
	case protocol.WorkspaceMessage_Remove:
		err = f.Remove(IOReqMsg.GetFileName())
	case protocol.WorkspaceMessage_IsExist:
		ok, err = f.IsExist(IOReqMsg.GetFileName())
	case protocol.WorkspaceMessage_Copy:
		err = f.Copy(IOReqMsg.GetCopyFileOldPath(), IOReqMsg.GetCopyFileNewPath())
	case protocol.WorkspaceMessage_IsDir:
		ok, err = f.IsDir(IOReqMsg.GetDirName())
	case protocol.WorkspaceMessage_Read:
		fileByte, err = f.Read(IOReqMsg.GetFileName())
	case protocol.WorkspaceMessage_ReadLines:
		fileByte, err = f.ReadLines(IOReqMsg.GetFileName(), IOReqMsg.GetReadLineBegin(), IOReqMsg.GetReadLineEnd())
	case protocol.WorkspaceMessage_WriteBytes:
		err = f.WriteBytes(IOReqMsg.GetFileName(), IOReqMsg.GetByteSlice())
	case protocol.WorkspaceMessage_AppendBytes:
		err = f.AppendBytes(IOReqMsg.GetFileName(), IOReqMsg.GetByteSlice())
	case protocol.WorkspaceMessage_WriteStrings:
		err = f.WriteStrings(IOReqMsg.GetFileName(), IOReqMsg.GetContent())
	case protocol.WorkspaceMessage_AppendStrings:
		err = f.AppendStrings(IOReqMsg.GetFileName(), IOReqMsg.GetContent())
	case protocol.WorkspaceMessage_Zip:
		err = f.Zip(IOReqMsg.GetZipName(), IOReqMsg.GetZipTargetFiles())
	case protocol.WorkspaceMessage_UnZip:
		err = f.UnZip(IOReqMsg.GetZipName(), IOReqMsg.GetZipTargetDir())
	case protocol.WorkspaceMessage_Gz:
		err = f.Gz(IOReqMsg.GetFileName())
	case protocol.WorkspaceMessage_UnGz:
		err = f.UnGz(IOReqMsg.GetFileName(), IOReqMsg.GetGzTargetFile())
	case protocol.WorkspaceMessage_Hash:
		fileByte, err = f.Hash(IOReqMsg.GetFileName())
	case protocol.WorkspaceMessage_CreateDir:
		err = f.MakeDir(IOReqMsg.GetDirName())
	case protocol.WorkspaceMessage_List:
		fileTree, err = f.List(IOReqMsg.GetDirName())
	}

	ioResponseMessage := &protocol.WorkspaceMessage_IOResponseMessage{
		Error:     message_utils.BuildErrorMessage(err),
		Operation: IOReqMsg.GetOperation(),
		Result:    ok,
		Data:      fileByte,
		FileTree:  fileTree,
	}
	w.distinct.Resource.Workspace = &protocol.WorkspaceMessage{IOResponse: ioResponseMessage}
}
