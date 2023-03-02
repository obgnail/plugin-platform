package release

import (
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/protocol"
	"github.com/obgnail/plugin-platform/common/utils/message"
	"github.com/obgnail/plugin-platform/host/resource/common"
)

var _ common_type.Workspace = (*Space)(nil)

type Space struct {
	plugin common_type.IPlugin
	sender common.Sender
}

func NewSpace(plugin common_type.IPlugin, sender common.Sender) common_type.Workspace {
	return &Space{plugin: plugin, sender: sender}
}

func (s *Space) buildMessage(ioRequest *protocol.WorkspaceMessage_IORequestMessage) *protocol.PlatformMessage {
	msg := message.GetInitMessage(nil, nil)
	msg.Resource = &protocol.ResourceMessage{
		Workspace: &protocol.WorkspaceMessage{IORequest: ioRequest},
	}
	return msg
}

func (s *Space) sendMsgToHost(platformMessage *protocol.PlatformMessage) (*protocol.PlatformMessage, common_type.PluginError) {
	return s.sender.Send(s.plugin, platformMessage)
}

func (s *Space) sendToHostAsync(platformMessage *protocol.PlatformMessage, callback common_type.AsyncInvokeCallbackParams) {
	cb := &spaceCallbackWrapper{Func: callback}
	s.sender.SendAsync(s.plugin, platformMessage, cb.callBack)
}

func (s *Space) CreateFile(name string) common_type.PluginError {
	ioRequestMessage := &protocol.WorkspaceMessage_IORequestMessage{
		Operation: protocol.WorkspaceMessage_Create,
		FileName:  name,
	}

	msg, err := s.sendMsgToHost(s.buildMessage(ioRequestMessage))
	if err != nil {
		return err
	}
	if retErr := msg.GetResource().GetWorkspace().GetIOResponse().GetError(); retErr != nil {
		return common_type.NewPluginError(int(retErr.Code), retErr.GetError(), retErr.GetMsg())
	}
	return nil
}

func (s *Space) MakeDir(dirname string) common_type.PluginError {
	ioRequestMessage := &protocol.WorkspaceMessage_IORequestMessage{
		Operation: protocol.WorkspaceMessage_CreateDir,
		DirName:   dirname,
	}

	msg, err := s.sendMsgToHost(s.buildMessage(ioRequestMessage))
	if err != nil {
		return err
	}
	if retErr := msg.GetResource().GetWorkspace().GetIOResponse().GetError(); retErr != nil {
		return common_type.NewPluginError(int(retErr.Code), retErr.GetError(), retErr.GetMsg())
	}
	return nil
}

func (s *Space) Rename(originalPath, newPath string) common_type.PluginError {
	ioRequestMessage := &protocol.WorkspaceMessage_IORequestMessage{
		Operation:   protocol.WorkspaceMessage_Rename,
		FileName:    originalPath,
		NewFileName: newPath,
	}
	msg, err := s.sendMsgToHost(s.buildMessage(ioRequestMessage))
	if err != nil {
		return err
	}
	if retErr := msg.GetResource().GetWorkspace().GetIOResponse().GetError(); retErr != nil {
		return common_type.NewPluginError(int(retErr.Code), retErr.GetError(), retErr.GetMsg())
	}
	return nil
}

func (s *Space) Remove(name string) common_type.PluginError {
	ioRequestMessage := &protocol.WorkspaceMessage_IORequestMessage{
		Operation: protocol.WorkspaceMessage_Remove,
		FileName:  name,
	}
	msg, err := s.sendMsgToHost(s.buildMessage(ioRequestMessage))
	if err != nil {
		return err
	}
	if retErr := msg.GetResource().GetWorkspace().GetIOResponse().GetError(); retErr != nil {
		return common_type.NewPluginError(int(retErr.Code), retErr.GetError(), retErr.GetMsg())
	}
	return nil
}

func (s *Space) IsExist(name string) (bool, common_type.PluginError) {
	ioRequestMessage := &protocol.WorkspaceMessage_IORequestMessage{
		Operation: protocol.WorkspaceMessage_IsExist,
		FileName:  name,
	}
	msg, err := s.sendMsgToHost(s.buildMessage(ioRequestMessage))
	if err != nil {
		return false, err
	}
	if retErr := msg.GetResource().GetWorkspace().GetIOResponse().GetError(); retErr != nil {
		return false, common_type.NewPluginError(int(retErr.Code), retErr.GetError(), retErr.GetMsg())
	}
	return msg.GetResource().GetWorkspace().GetIOResponse().GetResult(), nil
}

func (s *Space) IsDir(name string) (bool, common_type.PluginError) {
	ioRequestMessage := &protocol.WorkspaceMessage_IORequestMessage{
		Operation: protocol.WorkspaceMessage_IsDir,
		DirName:   name,
	}
	msg, err := s.sendMsgToHost(s.buildMessage(ioRequestMessage))
	if err != nil {
		return false, err
	}
	if retErr := msg.GetResource().GetWorkspace().GetIOResponse().GetError(); retErr != nil {
		return false, common_type.NewPluginError(int(retErr.Code), retErr.GetError(), retErr.GetMsg())
	}
	return msg.GetResource().GetWorkspace().GetIOResponse().GetResult(), nil
}

func (s *Space) Copy(originalPath string, newPath string) common_type.PluginError {
	ioRequestMessage := &protocol.WorkspaceMessage_IORequestMessage{
		Operation:       protocol.WorkspaceMessage_Copy,
		CopyFileOldPath: originalPath,
		CopyFileNewPath: newPath,
	}
	msg, err := s.sendMsgToHost(s.buildMessage(ioRequestMessage))
	if err != nil {
		return err
	}
	if retErr := msg.GetResource().GetWorkspace().GetIOResponse().GetError(); retErr != nil {
		return common_type.NewPluginError(int(retErr.Code), retErr.GetError(), retErr.GetMsg())
	}
	return nil
}

func (s *Space) Read(name string) ([]byte, common_type.PluginError) {
	ioRequestMessage := &protocol.WorkspaceMessage_IORequestMessage{
		Operation: protocol.WorkspaceMessage_Read,
		FileName:  name,
	}
	msg, err := s.sendMsgToHost(s.buildMessage(ioRequestMessage))
	if err != nil {
		return nil, err
	}
	if retErr := msg.GetResource().GetWorkspace().GetIOResponse().GetError(); retErr != nil {
		return nil, common_type.NewPluginError(int(retErr.Code), retErr.GetError(), retErr.GetMsg())
	}
	return msg.GetResource().GetWorkspace().GetIOResponse().GetData(), nil
}

func (s *Space) ReadLines(name string, readLineBegin, readLineEnd int32) ([]byte, common_type.PluginError) {
	ioRequestMessage := &protocol.WorkspaceMessage_IORequestMessage{
		Operation:     protocol.WorkspaceMessage_ReadLines,
		FileName:      name,
		ReadLineBegin: readLineBegin,
		ReadLineEnd:   readLineEnd,
	}
	msg, err := s.sendMsgToHost(s.buildMessage(ioRequestMessage))
	if err != nil {
		return nil, err
	}
	if retErr := msg.GetResource().GetWorkspace().GetIOResponse().GetError(); retErr != nil {
		return nil, common_type.NewPluginError(int(retErr.Code), retErr.GetError(), retErr.GetMsg())
	}
	return msg.GetResource().GetWorkspace().GetIOResponse().GetData(), nil
}

func (s *Space) WriteBytes(name string, byteSlice []byte) common_type.PluginError {
	ioRequestMessage := &protocol.WorkspaceMessage_IORequestMessage{
		Operation: protocol.WorkspaceMessage_WriteBytes,
		FileName:  name,
		ByteSlice: byteSlice,
	}
	msg, err := s.sendMsgToHost(s.buildMessage(ioRequestMessage))
	if err != nil {
		return err
	}
	if retErr := msg.GetResource().GetWorkspace().GetIOResponse().GetError(); retErr != nil {
		return common_type.NewPluginError(int(retErr.Code), retErr.GetError(), retErr.GetMsg())
	}
	return nil
}

func (s *Space) AppendBytes(name string, byteSlice []byte) common_type.PluginError {
	ioRequestMessage := &protocol.WorkspaceMessage_IORequestMessage{
		Operation: protocol.WorkspaceMessage_AppendBytes,
		FileName:  name,
		ByteSlice: byteSlice,
	}
	msg, err := s.sendMsgToHost(s.buildMessage(ioRequestMessage))
	if err != nil {
		return err
	}
	if retErr := msg.GetResource().GetWorkspace().GetIOResponse().GetError(); retErr != nil {
		return common_type.NewPluginError(int(retErr.Code), retErr.GetError(), retErr.GetMsg())
	}
	return nil
}

func (s *Space) WriteStrings(name string, content []string) common_type.PluginError {
	ioRequestMessage := &protocol.WorkspaceMessage_IORequestMessage{
		Operation: protocol.WorkspaceMessage_WriteStrings,
		FileName:  name,
		Content:   content,
	}
	msg, err := s.sendMsgToHost(s.buildMessage(ioRequestMessage))
	if err != nil {
		return err
	}
	if retErr := msg.GetResource().GetWorkspace().GetIOResponse().GetError(); retErr != nil {
		return common_type.NewPluginError(int(retErr.Code), retErr.GetError(), retErr.GetMsg())
	}
	return nil
}

func (s *Space) AppendStrings(name string, content []string) common_type.PluginError {
	ioRequestMessage := &protocol.WorkspaceMessage_IORequestMessage{
		Operation: protocol.WorkspaceMessage_AppendStrings,
		FileName:  name,
		Content:   content,
	}
	msg, err := s.sendMsgToHost(s.buildMessage(ioRequestMessage))
	if err != nil {
		return err
	}
	if retErr := msg.GetResource().GetWorkspace().GetIOResponse().GetError(); retErr != nil {
		return common_type.NewPluginError(int(retErr.Code), retErr.GetError(), retErr.GetMsg())
	}
	return nil
}

func (s *Space) Zip(outFileName string, targetFiles []string) common_type.PluginError {
	ioRequestMessage := &protocol.WorkspaceMessage_IORequestMessage{
		Operation:      protocol.WorkspaceMessage_Zip,
		ZipName:        outFileName,
		ZipTargetFiles: targetFiles,
	}
	msg, err := s.sendMsgToHost(s.buildMessage(ioRequestMessage))
	if err != nil {
		return err
	}
	if retErr := msg.GetResource().GetWorkspace().GetIOResponse().GetError(); retErr != nil {
		return common_type.NewPluginError(int(retErr.Code), retErr.GetError(), retErr.GetMsg())
	}
	return nil
}

func (s *Space) UnZip(name, targetDir string) common_type.PluginError {
	ioRequestMessage := &protocol.WorkspaceMessage_IORequestMessage{
		Operation:    protocol.WorkspaceMessage_UnZip,
		ZipName:      name,
		ZipTargetDir: targetDir,
	}
	msg, err := s.sendMsgToHost(s.buildMessage(ioRequestMessage))
	if err != nil {
		return err
	}
	if retErr := msg.GetResource().GetWorkspace().GetIOResponse().GetError(); retErr != nil {
		return common_type.NewPluginError(int(retErr.Code), retErr.GetError(), retErr.GetMsg())
	}
	return nil
}

func (s *Space) Gz(name string) common_type.PluginError {
	ioRequestMessage := &protocol.WorkspaceMessage_IORequestMessage{
		Operation: protocol.WorkspaceMessage_Gz,
		FileName:  name,
	}
	msg, err := s.sendMsgToHost(s.buildMessage(ioRequestMessage))
	if err != nil {
		return err
	}
	if retErr := msg.GetResource().GetWorkspace().GetIOResponse().GetError(); retErr != nil {
		return common_type.NewPluginError(int(retErr.Code), retErr.GetError(), retErr.GetMsg())
	}
	return nil
}

func (s *Space) UnGz(name, targetFile string) common_type.PluginError {
	ioRequestMessage := &protocol.WorkspaceMessage_IORequestMessage{
		Operation:    protocol.WorkspaceMessage_UnGz,
		FileName:     name,
		GzTargetFile: targetFile,
	}
	msg, err := s.sendMsgToHost(s.buildMessage(ioRequestMessage))
	if err != nil {
		return err
	}
	if retErr := msg.GetResource().GetWorkspace().GetIOResponse().GetError(); retErr != nil {
		return common_type.NewPluginError(int(retErr.Code), retErr.GetError(), retErr.GetMsg())
	}
	return nil
}

func (s *Space) Hash(name string) ([]byte, common_type.PluginError) {
	ioRequestMessage := &protocol.WorkspaceMessage_IORequestMessage{
		Operation: protocol.WorkspaceMessage_Hash,
		FileName:  name,
	}
	msg, err := s.sendMsgToHost(s.buildMessage(ioRequestMessage))
	if err != nil {
		return nil, err
	}
	if retErr := msg.GetResource().GetWorkspace().GetIOResponse().GetError(); retErr != nil {
		return nil, common_type.NewPluginError(int(retErr.Code), retErr.GetError(), retErr.GetMsg())
	}
	return msg.GetResource().GetWorkspace().GetIOResponse().GetData(), nil
}

func (s *Space) List(dirName string) ([]string, common_type.PluginError) {
	ioRequestMessage := &protocol.WorkspaceMessage_IORequestMessage{
		Operation: protocol.WorkspaceMessage_List,
		DirName:   dirName,
	}
	msg, err := s.sendMsgToHost(s.buildMessage(ioRequestMessage))
	if err != nil {
		return nil, err
	}
	if retErr := msg.GetResource().GetWorkspace().GetIOResponse().GetError(); retErr != nil {
		return nil, common_type.NewPluginError(int(retErr.Code), retErr.GetError(), retErr.GetMsg())
	}
	return msg.GetResource().GetWorkspace().GetIOResponse().GetFileTree(), nil
}

func (s *Space) AsyncCopy(originalPath string, newPath string, callback common_type.AsyncInvokeCallbackParams) {
	ioRequestMessage := &protocol.WorkspaceMessage_IORequestMessage{
		Operation:       protocol.WorkspaceMessage_Copy,
		CopyFileOldPath: originalPath,
		CopyFileNewPath: newPath,
	}
	s.sendToHostAsync(s.buildMessage(ioRequestMessage), callback)
}

func (s *Space) AsyncZip(outFileName string, targetFiles []string, callback common_type.AsyncInvokeCallbackParams) {
	ioRequestMessage := &protocol.WorkspaceMessage_IORequestMessage{
		Operation:      protocol.WorkspaceMessage_Zip,
		ZipName:        outFileName,
		ZipTargetFiles: targetFiles,
	}
	s.sendToHostAsync(s.buildMessage(ioRequestMessage), callback)
}

func (s *Space) AsyncUnZip(name, targetDir string, callback common_type.AsyncInvokeCallbackParams) {
	ioRequestMessage := &protocol.WorkspaceMessage_IORequestMessage{
		Operation:    protocol.WorkspaceMessage_UnZip,
		ZipName:      name,
		ZipTargetDir: targetDir,
	}
	s.sendToHostAsync(s.buildMessage(ioRequestMessage), callback)
}

func (s *Space) AsyncGz(name string, callback common_type.AsyncInvokeCallbackParams) {
	ioRequestMessage := &protocol.WorkspaceMessage_IORequestMessage{
		Operation: protocol.WorkspaceMessage_Gz,
		FileName:  name,
	}
	s.sendToHostAsync(s.buildMessage(ioRequestMessage), callback)
}

func (s *Space) AsyncUnGz(name, targetFile string, callback common_type.AsyncInvokeCallbackParams) {
	ioRequestMessage := &protocol.WorkspaceMessage_IORequestMessage{
		Operation:    protocol.WorkspaceMessage_UnGz,
		FileName:     name,
		GzTargetFile: targetFile,
	}
	s.sendToHostAsync(s.buildMessage(ioRequestMessage), callback)
}

type spaceCallbackWrapper struct {
	Func common_type.AsyncInvokeCallbackParams
}

func (w *spaceCallbackWrapper) callBack(input, result *protocol.PlatformMessage, err common_type.PluginError) {
	w.Func(err)
}
