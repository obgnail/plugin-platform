package local

import (
	"github.com/obgnail/plugin-platform/common/common_type"
)

var _ common_type.EventPublisher = (*Event)(nil)

type Event struct {
	plugin common_type.IPlugin
}

func NewEvent(plugin common_type.IPlugin) common_type.EventPublisher {
	return &Event{plugin: plugin}
}

func (event *Event) Subscribe(condition []string) common_type.PluginError {
	return nil
}

func (event *Event) SubscribeWithFilter(condition []string, filter map[string][]string) common_type.PluginError {
	return nil
}

func (event *Event) Unsubscribe(condition []string) common_type.PluginError {
	return nil
}
