package redis

import (
	"github.com/gomodule/redigo/redis"
	"github.com/obgnail/plugin-platform/common/errors"
)

func AddSet(conn redis.Conn, key string, members []string) error {
	l := len(members)
	if l > 0 {
		args := make([]interface{}, l+1)
		args[0] = key
		for i, s := range members {
			args[i+1] = s
		}
		conn.Send("SADD", args...)
	}
	conn.Flush()
	if _, err := conn.Receive(); err != nil {
		return errors.Trace(err)
	}
	return nil
}

func ListSetMembers(conn redis.Conn, key string) ([]string, error) {
	return redis.Strings(conn.Do("SMEMBERS", key))
}

func RemoveSetMembers(conn redis.Conn, key string, members []string) error {
	l := len(members)
	if l > 0 {
		args := make([]interface{}, l+1)
		args[0] = key
		for i, s := range members {
			args[i+1] = s
		}
		conn.Send("SREM", args...)
	}
	conn.Flush()
	if _, err := conn.Receive(); err != nil {
		return errors.Trace(err)
	}
	return nil
}
