package release

import (
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/message_utils"
	"github.com/obgnail/plugin-platform/common/protocol"
	"github.com/obgnail/plugin-platform/host/resource/common"
)

var _ common_type.EventPublisher = (*EventPublisher)(nil)

type EventPublisher struct {
	plugin common_type.IPlugin
	sender common.Sender
}

func NewEvent(plugin common_type.IPlugin, sender common.Sender) common_type.EventPublisher {
	return &EventPublisher{plugin: plugin, sender: sender}
}

func (event *EventPublisher) buildMessage(eventMsg *protocol.EventMessage) *protocol.PlatformMessage {
	msg := message_utils.GetInitMessage(nil, nil)
	msg.Resource = &protocol.ResourceMessage{Event: eventMsg}
	return msg
}

func (event *EventPublisher) sendMsgToHost(platformMessage *protocol.PlatformMessage) (*protocol.PlatformMessage, common_type.PluginError) {
	return event.sender.Send(event.plugin, platformMessage)
}

func (event *EventPublisher) send(eventMsg *protocol.EventMessage) common_type.PluginError {
	msg, err := event.sendMsgToHost(event.buildMessage(eventMsg))
	if err != nil {
		return err
	}
	retErr := msg.GetResource().GetEvent().GetError()
	if retErr != nil {
		return common_type.NewPluginError(int(retErr.Code), retErr.GetError(), retErr.GetMsg())
	}
	return nil
}

func (event *EventPublisher) Subscribe(condition []string) common_type.PluginError {
	eventMsg := &protocol.EventMessage{
		SubscribeOperation: protocol.EventMessage_Subscribe,
		Condition:          condition,
		SubscribeFilter:    nil,
	}
	return event.send(eventMsg)
}

func (event *EventPublisher) SubscribeWithFilter(condition []string, filter map[string][]string) common_type.PluginError {
	_filter := make(map[string]*protocol.Filter, len(filter))
	for key, value := range filter {
		_filter[key] = &protocol.Filter{Val: value}
	}

	eventMsg := &protocol.EventMessage{
		SubscribeOperation: protocol.EventMessage_SubscribeWithFilter,
		Condition:          condition,
		SubscribeFilter:    _filter,
	}
	return event.send(eventMsg)
}

func (event *EventPublisher) Unsubscribe(condition []string) common_type.PluginError {
	eventMsg := &protocol.EventMessage{
		SubscribeOperation: protocol.EventMessage_Unsubscribe,
		Condition:          condition,
		SubscribeFilter:    nil,
	}
	return event.send(eventMsg)
}
