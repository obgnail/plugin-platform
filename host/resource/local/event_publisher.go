package local

import (
	common "github.com/obgnail/plugin-platform/common_type"
)

var _ common.EventPublisher = (*Event)(nil)

type Event struct {
	plugin common.IPlugin
}

func NewEvent(plugin common.IPlugin) common.EventPublisher {
	return &Event{plugin: plugin}
}

func (event *Event) Subscribe(condition []string) common.PluginError {
	return nil
}

func (event *Event) SubscribeWithFilter(condition []string, filter map[string][]string) common.PluginError {
	return nil
}

func (event *Event) Unsubscribe(condition []string) common.PluginError {
	return nil
}
