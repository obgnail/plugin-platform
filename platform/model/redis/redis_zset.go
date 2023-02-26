package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/obgnail/plugin-platform/common/errors"
)

func GetScoreByMember(conn redis.Conn, key string, member string) (int64, error) {
	res, err := conn.Do("ZSCORE", key, member)
	if err != nil {
		return -1, errors.Trace(err)
	}

	var nilRes interface{}

	if res == nilRes {
		return 0, nil
	}
	score, err := redis.Int64(res, err)
	if err != nil {
		return -1, errors.Trace(err)
	}
	return score, nil
}

func MapScoreByMembersMinScore(conn redis.Conn, key string, minScore int64) (map[string]int64, error) {
	res, err := redis.Int64Map(conn.Do("ZRANGEBYSCORE", key, fmt.Sprintf("(%d", minScore), "+inf", "WITHSCORES"))
	if err != nil {
		return nil, errors.Trace(err)
	}
	return res, nil
}

func ZaddAndSetExpireTime(conn redis.Conn, key string, member string, score int64, expireTime int64) error {
	conn.Send("ZADD", key, score, member)
	conn.Send("EXPIRE", key, expireTime)
	conn.Flush()
	_, err := conn.Receive()
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

func Zadd(conn redis.Conn, key string, member string, score int64) error {
	conn.Send("ZADD", key, score, member)
	conn.Flush()
	_, err := conn.Receive()
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

func Publish(conn redis.Conn, channel string, msg string) {
	conn.Do("PUBLISH", channel, msg)
}

func Zincrby(conn redis.Conn, key string, member string, score int64) error {
	conn.Send("ZINCRBY", key, score, member)
	conn.Flush()
	_, err := conn.Receive()
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

func Zrem(conn redis.Conn, key string, member string) error {
	conn.Send("ZREM", key, member)
	conn.Flush()
	_, err := conn.Receive()
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

func GetByteByMember(conn redis.Conn, key string, minScore int64) ([][]byte, error) {
	res, err := redis.ByteSlices(conn.Do("ZRANGEBYSCORE", key, fmt.Sprintf("(%d", minScore), "+inf"))
	if err != nil {
		return nil, errors.Trace(err)
	}
	return res, nil
}

func GetByteByMemberRange(conn redis.Conn, key string, min, max int64) ([][]byte, error) {
	res, err := redis.ByteSlices(conn.Do("ZRANGE", key, min, max))
	if err != nil {
		return nil, errors.Trace(err)
	}
	return res, nil
}
