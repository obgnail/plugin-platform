package log

import (
	"github.com/obgnail/plugin-platform/common/errors"
	"github.com/obgnail/plugin-platform/common/log"
	"github.com/obgnail/plugin-platform/common/protocol"
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
	logger, err := NewLogger(l.appID, l.instanceID)
	if err != nil {
		log.ErrorDetails(errors.Trace(err))
		return
	}
	logMsg := l.Source.GetResource().GetLog()
	logger.Log(logMsg)
}
