package log

import (
	"fmt"
	"github.com/obgnail/plugin-platform/common/config"
	"github.com/obgnail/plugin-platform/common/errors"
	"github.com/obgnail/plugin-platform/common/log"
	"github.com/obgnail/plugin-platform/common/protocol"
	"github.com/obgnail/plugin-platform/common/utils/file_path"
	"os"
)

type Log struct {
	Source     *protocol.PlatformMessage
	Distinct   *protocol.PlatformMessage
	appID      string
	instanceID string
}

func NewLog(sourceMessage *protocol.PlatformMessage, distinctMessage *protocol.PlatformMessage) *Log {
	dataBase := &Log{
		Source:     sourceMessage,
		Distinct:   distinctMessage,
		appID:      sourceMessage.GetResource().GetSender().GetApplication().GetApplicationID(),
		instanceID: sourceMessage.GetResource().GetSender().GetInstanceID(),
	}
	return dataBase
}

func (l *Log) Execute() {
	path, err := l.getPath()
	if err != nil {
		log.ErrorDetails(errors.Trace(err))
		return
	}

	logger, err := NewLogger(path)
	if err != nil {
		log.ErrorDetails(errors.Trace(err))
		return
	}
	logMsg := l.Source.GetResource().GetLog()
	logger.Log(logMsg)
}

func (l *Log) getPath() (string, error) {
	dirPath := config.StringOrPanic("platform.plugin_log_dir")
	dirPath = file_path.JoinPath(dirPath, l.appID)
	if err := os.MkdirAll(dirPath, 0640); err != nil {
		return "", errors.Trace(err)
	}
	file := file_path.JoinPath(dirPath, fmt.Sprintf("%s.log", l.instanceID))
	return file, nil
}
