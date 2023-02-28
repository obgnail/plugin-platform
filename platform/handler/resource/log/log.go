package log

import (
	"github.com/obgnail/plugin-platform/common/errors"
	commonLog "github.com/obgnail/plugin-platform/common/log"
	"github.com/obgnail/plugin-platform/common/protocol"
	"io"
	"os"
)

type Logger struct {
	*commonLog.Logger
}

func NewLogger(file string) (*Logger, error) {
	sep := "/github.com"
	pathPrefix := "/github.com/obgnail/plugin-platform"
	logFile, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0640)
	if err != nil {
		return nil, errors.Trace(err)
	}

	auditConfig := commonLog.Config{
		Sep:        sep,
		Level:      commonLog.TraceLevel,
		Target:     io.MultiWriter(os.Stdout, logFile),
		PathPrefix: pathPrefix,
		Encoder:    &commonLog.JsonEncoder{EnableBuffer: false}, //先禁用buffer,如果开启需要处理系统信号量
	}

	l, _, err := commonLog.NewLogger(auditConfig)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return &Logger{Logger: l}, nil
}

func (l *Logger) Log(msg *protocol.LogMessage) {
	level := msg.GetLevel()
	content := msg.GetContent()

	switch level {
	case protocol.LogMessage_Trace:
		l.Trace(content)
	case protocol.LogMessage_Info:
		l.Info(content)
	case protocol.LogMessage_Warning:
		l.Warn(content)
	case protocol.LogMessage_Error:
		l.Error(content)
	}
}
