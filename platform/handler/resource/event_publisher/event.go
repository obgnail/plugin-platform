package event_publisher

import (
	"fmt"
	"github.com/obgnail/plugin-platform/common/errors"
	"github.com/obgnail/plugin-platform/common/math"
	"github.com/obgnail/plugin-platform/platform/model/redis"
)

const (
	eventTopicKey = "event_topic" // map[event_topic][]string{event_subscriber}
)

var Events = []string{
	"project.task",
	"project.user",
	"project.issue",
	"project.article",
}

// Subscribe 不支持 wildcard
func Subscribe(condition []string, instanceID string) error {
	var cnds []string
	for _, cnd := range condition {
		for _, event := range Events {
			if event == cnd {
				cnds = append(cnds, cnd)
			}
		}
	}
	if err := onSubscribe(cnds, instanceID); err != nil {
		return errors.Trace(err)
	}
	return nil
}

// Unsubscribe 支持 wildcard
func Unsubscribe(condition []string, instanceID string) error {
	var cnds []string
	for _, cnd := range condition {
		for _, event := range Events {
			if ok := math.MatchWildcard(event, cnd); ok {
				cnds = append(cnds, event)
			}
		}
	}
	if err := onUnsubscribe(cnds, instanceID); err != nil {
		return errors.Trace(err)
	}
	return nil
}

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

func onPublish(topic string) ([]string, error) {
	conn := redis.GetRedisConn()
	defer conn.Close()

	key := fmt.Sprintf("%s:%s", eventTopicKey, topic)
	result, err := redis.ListSetMembers(conn, key)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return result, nil
}
