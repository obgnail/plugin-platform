package event_publisher

import (
	"fmt"
	"github.com/obgnail/plugin-platform/common/errors"
	"github.com/obgnail/plugin-platform/common/log"
	"github.com/obgnail/plugin-platform/common/utils/math"
	"github.com/obgnail/plugin-platform/platform/model/mysql"
	"github.com/obgnail/plugin-platform/platform/model/redis"
)

const (
	eventTopicKey = "event_topic" // event_topic:【topic】:【instanceID】
)

// Subscribe 不支持 wildcard
func Subscribe(condition []string, instanceID string) error {
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
	if err := onSubscribe(cnds, instanceID); err != nil {
		return errors.Trace(err)
	}
	return nil
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
	if err := onUnsubscribe(cnds, instanceID); err != nil {
		return errors.Trace(err)
	}
	return nil
}

// 写入 redis
func onSubscribe(topics []string, subscribe string) error {
	conn := redis.GetRedisConn()
	defer conn.Close()

	for _, topic := range topics {
		key := fmt.Sprintf("%s:%s", eventTopicKey, topic)
		err := redis.AddSet(conn, key, []string{subscribe})
		if err != nil {
			return errors.Trace(err)
		}
	}

	return nil
}

func onUnsubscribe(topics []string, unsubscribe string) error {
	conn := redis.GetRedisConn()
	defer conn.Close()

	for _, topic := range topics {
		key := fmt.Sprintf("%s:%s", eventTopicKey, topic)
		err := redis.RemoveSetMembers(conn, key, []string{unsubscribe})
		if err != nil {
			return errors.Trace(err)
		}
	}

	return nil
}

func ListInstances(topic string) ([]string, error) {
	conn := redis.GetRedisConn()
	defer conn.Close()

	key := fmt.Sprintf("%s:%s", eventTopicKey, topic)
	result, err := redis.ListSetMembers(conn, key)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return result, nil
}

func Filter(topic string, instanceID string) (bool, error) {
	instances, err := ListInstances(topic)
	if err != nil {
		return false, errors.Trace(err)
	}
	for _, id := range instances {
		if instanceID == id {
			return true, nil
		}
	}
	return false, nil
}
