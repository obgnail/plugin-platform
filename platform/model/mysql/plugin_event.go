package mysql

import "github.com/obgnail/plugin-platform/common/errors"

func ModelPluginEvent() *PluginEvent {
	var e = new(PluginEvent)
	e.Child = e
	return e
}

type PluginEvent struct {
	BaseModel
	Action   string `gorm:"action" json:"action"`
	Extended string `gorm:"extended" json:"extended"`
}

func (e *PluginEvent) tableName() string {
	return "plugin_event"
}

func GetEvents() ([]*PluginEvent, error) {
	var events = make([]*PluginEvent, 0)
	err := ModelPluginEvent().All(&events, &PluginEvent{})
	if err != nil {
		return nil, errors.Trace(err)
	}
	return events, nil
}
