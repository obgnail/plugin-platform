package event

import (
	"fmt"
	"github.com/obgnail/plugin-platform/common/errors"
	"github.com/obgnail/plugin-platform/common/log"
	"github.com/obgnail/plugin-platform/common/utils/math"
	"github.com/obgnail/plugin-platform/platform/model/mysql"
	"strings"
	"sync"
)

var _event *event

type event struct {
	// 【topic_instanceID】:【filter】
	// map[string]map[string][]string
	m sync.Map

	// filter的逻辑与主系统高度绑定。需要用户自行定义。
	// 例如当topic是project.task, payload是task有关信息, filter是task_uuid_in:["ABCDE"] 该怎么处理
	// 这种处理需要了解payload的数据结构,因此抽象出Filter函数供用户自定义。
	f Filter
}

func InitEvent() {
	_event = &event{f: DummyFilter}
}

func getKey(topic string, subscribe string) string {
	return fmt.Sprintf("%s_%s", topic, subscribe)
}

func splitKey(s string) (string, string) {
	idx := strings.LastIndex(s, "_")
	return s[:idx], s[idx+1:]
}

// Subscribe 不支持 wildcard
func Subscribe(condition []string, instanceID string, filter map[string][]string) error {
	events, err := mysql.GetEvents()
	if err != nil {
		return errors.Trace(err)
	}

	var cnds []string
	for _, cnd := range condition {
		for _, event := range events {
			if event.Action == cnd {
				cnds = append(cnds, cnd)
			}
		}
	}

	log.Trace("[%s] Subscribe:[%v]", instanceID, condition)

	for _, cnd := range cnds {
		onSubscribe(cnd, instanceID, filter)
	}
	return nil
}

func onSubscribe(topic, instanceID string, filter map[string][]string) {
	key := getKey(topic, instanceID)

	f, ok := _event.m.Load(key)
	if !ok {
		_event.m.Store(key, filter)
		return
	}

	_filter := f.(map[string][]string)
	for k, v := range filter {
		_filter[k] = v
	}
	_event.m.Store(key, _filter)
}

// Unsubscribe 支持 wildcard
func Unsubscribe(condition []string, instanceID string) error {
	events, err := mysql.GetEvents()
	if err != nil {
		return errors.Trace(err)
	}

	var cnds []string
	for _, cnd := range condition {
		for _, event := range events {
			if ok := math.MatchWildcard(event.Action, cnd); ok {
				cnds = append(cnds, event.Action)
			}
		}
	}
	log.Trace("[%s] Unsubscribe:[%v]", instanceID, condition)
	for _, cnd := range cnds {
		key := getKey(cnd, instanceID)
		_event.m.Delete(key)
	}
	return nil
}

func FilterEvent(instanceID, topic string, payload []byte) bool {
	key := getKey(topic, instanceID)
	filter, ok := _event.m.Load(key)
	if !ok {
		return false
	}
	if filter == nil {
		return true
	}

	_filter := filter.(map[string][]string)
	return _event.f(topic, payload, _filter)
}
