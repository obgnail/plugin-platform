package release

import (
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/protocol"
	"github.com/obgnail/plugin-platform/common/utils/message"
	"github.com/obgnail/plugin-platform/host/resource/common"
)

var _ common_type.PluginLogger = (*Logger)(nil)

type Logger struct {
	plugin common_type.IPlugin
	sender common.Sender
}

func NewLogger(plugin common_type.IPlugin, sender common.Sender) *Logger {
	l := &Logger{plugin: plugin, sender: sender}
	return l
}

func (l *Logger) buildMessage(logMsg *protocol.LogMessage) *protocol.PlatformMessage {
	msg := message.GetInitMessage(nil, nil)
	msg.Resource = &protocol.ResourceMessage{
		Log: logMsg,
	}
	return msg
}

func (l *Logger) sendMsgToHost(platformMessage *protocol.PlatformMessage) (*protocol.PlatformMessage, common_type.PluginError) {
	return l.sender.Send(l.plugin, platformMessage)
}

func (l *Logger) send(level protocol.LogMessage_LogLevel, content string) {
	logMsg := &protocol.LogMessage{
		Level:   level,
		Content: content,
	}
	l.sendMsgToHost(l.buildMessage(logMsg))
}

func (l *Logger) Trace(content string) {
	l.send(protocol.LogMessage_Trace, content)
}

func (l *Logger) Info(content string) {
	l.send(protocol.LogMessage_Info, content)
}

func (l *Logger) Warn(content string) {
	l.send(protocol.LogMessage_Warning, content)
}

func (l *Logger) Error(content string) {
	l.send(protocol.LogMessage_Error, content)
}
