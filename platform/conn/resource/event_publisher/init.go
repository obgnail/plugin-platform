package event_publisher

import (
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/errors"
	"github.com/obgnail/plugin-platform/common/log"
	"github.com/obgnail/plugin-platform/common/protocol"
	"github.com/obgnail/plugin-platform/common/utils/message"
)

type Event struct {
	source     *protocol.PlatformMessage
	distinct   *protocol.PlatformMessage
	opType     protocol.EventMessage_SubscribeOperationType
	AppID      string
	instanceID string
}

func NewEvent(source, distinct *protocol.PlatformMessage) *Event {
	e := &Event{
		source:     source,
		distinct:   distinct,
		opType:     source.GetResource().GetEvent().GetSubscribeOperation(),
		AppID:      source.GetResource().GetSender().GetApplication().GetApplicationID(),
		instanceID: source.GetResource().GetSender().GetInstanceID(),
	}
	return e
}

func (e *Event) Execute() {
	var err error
	switch e.opType {
	case protocol.EventMessage_Subscribe:
		err = e.Subscribe()
	case protocol.EventMessage_Unsubscribe:
		err = e.Unsubscribe()
	case protocol.EventMessage_SubscribeWithFilter:
		err = e.SubscribeWithFilter()
	}
	e.buildMsg(err)
}

func (e *Event) Subscribe() error {
	conditions := e.source.GetResource().GetEvent().GetCondition()
	if err := Subscribe(conditions, e.instanceID); err != nil {
		return errors.Trace(err)
	}
	return nil
}

func (e *Event) Unsubscribe() error {
	conditions := e.source.GetResource().GetEvent().GetCondition()
	if err := Unsubscribe(conditions, e.instanceID); err != nil {
		return errors.Trace(err)
	}
	return nil
}

func (e *Event) SubscribeWithFilter() error {
	conditions := e.source.GetResource().GetEvent().GetCondition()
	if err := Subscribe(conditions, e.instanceID); err != nil {
		return errors.Trace(err)
	}
	return nil
}

func (e *Event) buildMsg(err error) {
	msg := &protocol.EventMessage{
		SubscribeOperation: e.opType,
		Condition:          e.source.GetResource().GetEvent().GetCondition(),
		SubscribeFilter:    e.source.GetResource().GetEvent().GetSubscribeFilter(),
	}
	if err != nil {
		log.ErrorDetails(errors.Trace(err))
		e := common_type.NewPluginError(common_type.NotifyEventFailure, err.Error())
		msg.Error = message.BuildErrorMessage(e)
	}
	e.distinct.Resource.Event = msg
}
