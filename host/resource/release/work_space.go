package release

import (
	common "github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/protocol"
	"github.com/obgnail/plugin-platform/host/handler"
	"github.com/obgnail/plugin-platform/host/resource/utils"
)

type Space struct {
	msg     *protocol.WorkspaceMessage_IORequestMessage
	plugin  common.IPlugin
	handler *handler.HostHandler
}

func NewSpace(plugin common.IPlugin, handler *handler.HostHandler) common.Workspace {
	return &Space{plugin: plugin, handler: handler}
}

func (f *Space) buildMessage(ioRequest *protocol.WorkspaceMessage_IORequestMessage) *protocol.PlatformMessage {
	msg := utils.GetHostToPlatFormMessage()
	msg.Resource = &protocol.ResourceMessage{
		Workspace: &protocol.WorkspaceMessage{
			IORequest: ioRequest,
		},
	}
	return msg
}

func (f *Space) sendMsgToHost(platformMessage *protocol.PlatformMessage) (*protocol.PlatformMessage, common.PluginError) {
	return f.handler.Send(f.plugin, platformMessage)
}

func (f *Space) sendToHostAsync(platformMessage *protocol.PlatformMessage, callback common.AsyncInvokeCallbackParams) {
	cb := &callbackWrapper{Func: callback}
	f.handler.SendAsync(f.plugin, platformMessage, cb.callBack)
}

func (f *Space) CreateFile(name string) common.PluginError {
	ioRequestMessage := &protocol.WorkspaceMessage_IORequestMessage{
		Operation: protocol.WorkspaceMessage_Create,
		FileName:  name,
	}

	msg, err := f.sendMsgToHost(f.buildMessage(ioRequestMessage))
	if err != nil {
		return err
	}
	if retErr := msg.GetResource().GetWorkspace().GetIOResponse().GetError(); retErr != nil {
		return common.NewPluginError(int(retErr.Code), retErr.GetError(), retErr.GetMsg())
	}
	return nil
}

func (f *Space) MakeDir(dirname string) common.PluginError {
	ioRequestMessage := &protocol.WorkspaceMessage_IORequestMessage{
		Operation: protocol.WorkspaceMessage_CreateDir,
		DirName:   dirname,
	}

	msg, err := f.sendMsgToHost(f.buildMessage(ioRequestMessage))
	if err != nil {
		return err
	}
	if retErr := msg.GetResource().GetWorkspace().GetIOResponse().GetError(); retErr != nil {
		return common.NewPluginError(int(retErr.Code), retErr.GetError(), retErr.GetMsg())
	}
	return nil
}

func (f *Space) Rename(originalPath, newPath string) common.PluginError {
	ioRequestMessage := &protocol.WorkspaceMessage_IORequestMessage{
		Operation:   protocol.WorkspaceMessage_Rename,
		FileName:    originalPath,
		NewFileName: newPath,
	}
	msg, err := f.sendMsgToHost(f.buildMessage(ioRequestMessage))
	if err != nil {
		return err
	}
	if retErr := msg.GetResource().GetWorkspace().GetIOResponse().GetError(); retErr != nil {
		return common.NewPluginError(int(retErr.Code), retErr.GetError(), retErr.GetMsg())
	}
	return nil
}

func (f *Space) Remove(name string) common.PluginError {
	ioRequestMessage := &protocol.WorkspaceMessage_IORequestMessage{
		Operation: protocol.WorkspaceMessage_Remove,
		FileName:  name,
	}
	msg, err := f.sendMsgToHost(f.buildMessage(ioRequestMessage))
	if err != nil {
		return err
	}
	if retErr := msg.GetResource().GetWorkspace().GetIOResponse().GetError(); retErr != nil {
		return common.NewPluginError(int(retErr.Code), retErr.GetError(), retErr.GetMsg())
	}
	return nil
}

func (f *Space) IsExist(name string) (bool, common.PluginError) {
	ioRequestMessage := &protocol.WorkspaceMessage_IORequestMessage{
		Operation: protocol.WorkspaceMessage_IsExist,
		FileName:  name,
	}
	msg, err := f.sendMsgToHost(f.buildMessage(ioRequestMessage))
	if err != nil {
		return false, err
	}
	if retErr := msg.GetResource().GetWorkspace().GetIOResponse().GetError(); retErr != nil {
		return false, common.NewPluginError(int(retErr.Code), retErr.GetError(), retErr.GetMsg())
	}
	return msg.GetResource().GetWorkspace().GetIOResponse().GetResult(), nil
}

func (f *Space) IsDir(name string) (bool, common.PluginError) {
	ioRequestMessage := &protocol.WorkspaceMessage_IORequestMessage{
		Operation: protocol.WorkspaceMessage_IsDir,
		DirName:   name,
	}
	msg, err := f.sendMsgToHost(f.buildMessage(ioRequestMessage))
	if err != nil {
		return false, err
	}
	if retErr := msg.GetResource().GetWorkspace().GetIOResponse().GetError(); retErr != nil {
		return false, common.NewPluginError(int(retErr.Code), retErr.GetError(), retErr.GetMsg())
	}
	return msg.GetResource().GetWorkspace().GetIOResponse().GetResult(), nil
}

func (f *Space) Copy(originalPath string, newPath string) common.PluginError {
	ioRequestMessage := &protocol.WorkspaceMessage_IORequestMessage{
		Operation:       protocol.WorkspaceMessage_Copy,
		CopyFileOldPath: originalPath,
		CopyFileNewPath: newPath,
	}
	msg, err := f.sendMsgToHost(f.buildMessage(ioRequestMessage))
	if err != nil {
		return err
	}
	if retErr := msg.GetResource().GetWorkspace().GetIOResponse().GetError(); retErr != nil {
		return common.NewPluginError(int(retErr.Code), retErr.GetError(), retErr.GetMsg())
	}
	return nil
}

func (f *Space) Read(name string) ([]byte, common.PluginError) {
	ioRequestMessage := &protocol.WorkspaceMessage_IORequestMessage{
		Operation: protocol.WorkspaceMessage_Read,
		FileName:  name,
	}
	msg, err := f.sendMsgToHost(f.buildMessage(ioRequestMessage))
	if err != nil {
		return nil, err
	}
	if retErr := msg.GetResource().GetWorkspace().GetIOResponse().GetError(); retErr != nil {
		return nil, common.NewPluginError(int(retErr.Code), retErr.GetError(), retErr.GetMsg())
	}
	return msg.GetResource().GetWorkspace().GetIOResponse().GetData(), nil
}

func (f *Space) ReadLines(name string, readLineBegin, readLineEnd int32) ([]byte, common.PluginError) {
	ioRequestMessage := &protocol.WorkspaceMessage_IORequestMessage{
		Operation:     protocol.WorkspaceMessage_ReadLines,
		FileName:      name,
		ReadLineBegin: readLineBegin,
		ReadLineEnd:   readLineEnd,
	}
	msg, err := f.sendMsgToHost(f.buildMessage(ioRequestMessage))
	if err != nil {
		return nil, err
	}
	if retErr := msg.GetResource().GetWorkspace().GetIOResponse().GetError(); retErr != nil {
		return nil, common.NewPluginError(int(retErr.Code), retErr.GetError(), retErr.GetMsg())
	}
	return msg.GetResource().GetWorkspace().GetIOResponse().GetData(), nil
}

func (f *Space) WriteBytes(name string, byteSlice []byte) common.PluginError {
	ioRequestMessage := &protocol.WorkspaceMessage_IORequestMessage{
		Operation: protocol.WorkspaceMessage_WriteBytes,
		FileName:  name,
		ByteSlice: byteSlice,
	}
	msg, err := f.sendMsgToHost(f.buildMessage(ioRequestMessage))
	if err != nil {
		return err
	}
	if retErr := msg.GetResource().GetWorkspace().GetIOResponse().GetError(); retErr != nil {
		return common.NewPluginError(int(retErr.Code), retErr.GetError(), retErr.GetMsg())
	}
	return nil
}

func (f *Space) AppendBytes(name string, byteSlice []byte) common.PluginError {
	ioRequestMessage := &protocol.WorkspaceMessage_IORequestMessage{
		Operation: protocol.WorkspaceMessage_AppendBytes,
		FileName:  name,
		ByteSlice: byteSlice,
	}
	msg, err := f.sendMsgToHost(f.buildMessage(ioRequestMessage))
	if err != nil {
		return err
	}
	if retErr := msg.GetResource().GetWorkspace().GetIOResponse().GetError(); retErr != nil {
		return common.NewPluginError(int(retErr.Code), retErr.GetError(), retErr.GetMsg())
	}
	return nil
}

func (f *Space) WriteStrings(name string, content []string) common.PluginError {
	ioRequestMessage := &protocol.WorkspaceMessage_IORequestMessage{
		Operation: protocol.WorkspaceMessage_WriteStrings,
		FileName:  name,
		Content:   content,
	}
	msg, err := f.sendMsgToHost(f.buildMessage(ioRequestMessage))
	if err != nil {
		return err
	}
	if retErr := msg.GetResource().GetWorkspace().GetIOResponse().GetError(); retErr != nil {
		return common.NewPluginError(int(retErr.Code), retErr.GetError(), retErr.GetMsg())
	}
	return nil
}

func (f *Space) AppendStrings(name string, content []string) common.PluginError {
	ioRequestMessage := &protocol.WorkspaceMessage_IORequestMessage{
		Operation: protocol.WorkspaceMessage_AppendStrings,
		FileName:  name,
		Content:   content,
	}
	msg, err := f.sendMsgToHost(f.buildMessage(ioRequestMessage))
	if err != nil {
		return err
	}
	if retErr := msg.GetResource().GetWorkspace().GetIOResponse().GetError(); retErr != nil {
		return common.NewPluginError(int(retErr.Code), retErr.GetError(), retErr.GetMsg())
	}
	return nil
}

func (f *Space) Zip(outFileName string, targetFiles []string) common.PluginError {
	ioRequestMessage := &protocol.WorkspaceMessage_IORequestMessage{
		Operation:      protocol.WorkspaceMessage_Zip,
		ZipName:        outFileName,
		ZipTargetFiles: targetFiles,
	}
	msg, err := f.sendMsgToHost(f.buildMessage(ioRequestMessage))
	if err != nil {
		return err
	}
	if retErr := msg.GetResource().GetWorkspace().GetIOResponse().GetError(); retErr != nil {
		return common.NewPluginError(int(retErr.Code), retErr.GetError(), retErr.GetMsg())
	}
	return nil
}

func (f *Space) UnZip(name, targetDir string) common.PluginError {
	ioRequestMessage := &protocol.WorkspaceMessage_IORequestMessage{
		Operation:    protocol.WorkspaceMessage_UnZip,
		ZipName:      name,
		ZipTargetDir: targetDir,
	}
	msg, err := f.sendMsgToHost(f.buildMessage(ioRequestMessage))
	if err != nil {
		return err
	}
	if retErr := msg.GetResource().GetWorkspace().GetIOResponse().GetError(); retErr != nil {
		return common.NewPluginError(int(retErr.Code), retErr.GetError(), retErr.GetMsg())
	}
	return nil
}

func (f *Space) Gz(name string) common.PluginError {
	ioRequestMessage := &protocol.WorkspaceMessage_IORequestMessage{
		Operation: protocol.WorkspaceMessage_Gz,
		FileName:  name,
	}
	msg, err := f.sendMsgToHost(f.buildMessage(ioRequestMessage))
	if err != nil {
		return err
	}
	if retErr := msg.GetResource().GetWorkspace().GetIOResponse().GetError(); retErr != nil {
		return common.NewPluginError(int(retErr.Code), retErr.GetError(), retErr.GetMsg())
	}
	return nil
}

func (f *Space) UnGz(name, targetFile string) common.PluginError {
	ioRequestMessage := &protocol.WorkspaceMessage_IORequestMessage{
		Operation:    protocol.WorkspaceMessage_UnGz,
		FileName:     name,
		GzTargetFile: targetFile,
	}
	msg, err := f.sendMsgToHost(f.buildMessage(ioRequestMessage))
	if err != nil {
		return err
	}
	if retErr := msg.GetResource().GetWorkspace().GetIOResponse().GetError(); retErr != nil {
		return common.NewPluginError(int(retErr.Code), retErr.GetError(), retErr.GetMsg())
	}
	return nil
}

func (f *Space) Hash(name string) ([]byte, common.PluginError) {
	ioRequestMessage := &protocol.WorkspaceMessage_IORequestMessage{
		Operation: protocol.WorkspaceMessage_Hash,
		FileName:  name,
	}
	msg, err := f.sendMsgToHost(f.buildMessage(ioRequestMessage))
	if err != nil {
		return nil, err
	}
	if retErr := msg.GetResource().GetWorkspace().GetIOResponse().GetError(); retErr != nil {
		return nil, common.NewPluginError(int(retErr.Code), retErr.GetError(), retErr.GetMsg())
	}
	return msg.GetResource().GetWorkspace().GetIOResponse().GetData(), nil
}

func (f *Space) List(dirName string) ([]string, common.PluginError) {
	ioRequestMessage := &protocol.WorkspaceMessage_IORequestMessage{
		Operation: protocol.WorkspaceMessage_List,
		DirName:   dirName,
	}
	msg, err := f.sendMsgToHost(f.buildMessage(ioRequestMessage))
	if err != nil {
		return nil, err
	}
	if retErr := msg.GetResource().GetWorkspace().GetIOResponse().GetError(); retErr != nil {
		return nil, common.NewPluginError(int(retErr.Code), retErr.GetError(), retErr.GetMsg())
	}
	return msg.GetResource().GetWorkspace().GetIOResponse().GetFileTree(), nil
}

func (f *Space) AsyncCopy(originalPath string, newPath string, callback common.AsyncInvokeCallbackParams) {
	ioRequestMessage := &protocol.WorkspaceMessage_IORequestMessage{
		Operation:       protocol.WorkspaceMessage_Copy,
		CopyFileOldPath: originalPath,
		CopyFileNewPath: newPath,
	}
	f.sendToHostAsync(f.buildMessage(ioRequestMessage), callback)
}

func (f *Space) AsyncZip(outFileName string, targetFiles []string, callback common.AsyncInvokeCallbackParams) {
	ioRequestMessage := &protocol.WorkspaceMessage_IORequestMessage{
		Operation:      protocol.WorkspaceMessage_Zip,
		ZipName:        outFileName,
		ZipTargetFiles: targetFiles,
	}
	f.sendToHostAsync(f.buildMessage(ioRequestMessage), callback)
}

func (f *Space) AsyncUnZip(name, targetDir string, callback common.AsyncInvokeCallbackParams) {
	ioRequestMessage := &protocol.WorkspaceMessage_IORequestMessage{
		Operation:    protocol.WorkspaceMessage_UnZip,
		ZipName:      name,
		ZipTargetDir: targetDir,
	}
	f.sendToHostAsync(f.buildMessage(ioRequestMessage), callback)
}

func (f *Space) AsyncGz(name string, callback common.AsyncInvokeCallbackParams) {
	ioRequestMessage := &protocol.WorkspaceMessage_IORequestMessage{
		Operation: protocol.WorkspaceMessage_Gz,
		FileName:  name,
	}
	f.sendToHostAsync(f.buildMessage(ioRequestMessage), callback)
}

func (f *Space) AsyncUnGz(name, targetFile string, callback common.AsyncInvokeCallbackParams) {
	ioRequestMessage := &protocol.WorkspaceMessage_IORequestMessage{
		Operation:    protocol.WorkspaceMessage_UnGz,
		FileName:     name,
		GzTargetFile: targetFile,
	}
	f.sendToHostAsync(f.buildMessage(ioRequestMessage), callback)
}

type callbackWrapper struct {
	Func common.AsyncInvokeCallbackParams
}

func (w *callbackWrapper) callBack(input, result *protocol.PlatformMessage, err common.PluginError) {
	w.Func(err)
}
