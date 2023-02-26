package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/obgnail/plugin-platform/common/errors"
)

func DeleteHashKeyFields(conn redis.Conn, args []interface{}) error {
	conn.Send("HDEL", args...)
	conn.Flush()
	_, err := conn.Receive()
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

func GetHashKey(conn redis.Conn, key string) ([]interface{}, error) {
	values, err := redis.Values(conn.Do("HGETALL", key))
	if err != nil {
		return nil, errors.Trace(err)
	}
	return values, nil
}

func ListHashKey(conn redis.Conn, key string, fields []interface{}) ([]string, error) {
	args := redis.Args{}
	args = args.Add(key)
	args = args.Add(fields...)
	values, err := redis.Strings(conn.Do("HMGET", args...))
	if err != nil {
		return nil, errors.Trace(err)
	}
	return values, nil
}

func GetHashIntValue(conn redis.Conn, key, field string) (bool, int, error) {
	value, err := redis.Int(conn.Do("HGET", key, field))
	if err == redis.ErrNil {
		return false, -1, nil
	}
	if err != nil {
		return false, -1, errors.Trace(err)
	}
	return true, value, nil
}

func GetHashStringValue(conn redis.Conn, key, field string) (bool, string, error) {
	value, err := redis.String(conn.Do("HGET", key, field))
	if err == redis.ErrNil {
		return false, "", nil
	}
	if err != nil {
		return false, "", errors.Trace(err)
	}
	return true, value, nil
}

func AddHashKeys(conn redis.Conn, expires bool, records [][]interface{}) error {
	length := len(records)
	for _, record := range records {
		if len(record) == 0 {
			return errors.Errorf(errors.RedisError, "invalid args")
		}
		conn.Send("HSET", record...)
		if expires {
			conn.Send("EXPIRE", record[0], redisKeyExpire)
		}
	}
	conn.Flush()
	for i := 0; i < length; i++ {
		_, err := conn.Receive()
		if err != nil {
			return errors.Trace(err)
		}
	}
	return nil
}

func AddHashKeysAndSetExpireTime(conn redis.Conn, expireTime int64, records [][]interface{}) error {
	length := len(records)
	if length == 0 {
		return nil
	}
	for _, record := range records {
		if len(record) == 0 {
			return errors.Errorf(errors.RedisError, "invalid args")
		}
		conn.Send("HSET", record...)
		conn.Send("EXPIRE", record[0], expireTime)
	}
	conn.Flush()
	for i := 0; i < length; i++ {
		_, err := conn.Receive()
		if err != nil {
			return errors.Trace(err)
		}
	}
	return nil
}

func BatchAddHashKeysAndSetExpireTime(conn redis.Conn, expireTime int64, key string, fieldAndValue []interface{}) error {
	length := len(fieldAndValue)
	if length == 0 {
		return nil
	}
	args := redis.Args{}
	args = args.Add(key)
	args = args.Add(fieldAndValue...)
	conn.Send("HMSET", args...)
	conn.Send("EXPIRE", key, expireTime)
	conn.Flush()
	_, err := conn.Receive()
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

func DeleteKey(conn redis.Conn, key string) error {
	if _, err := conn.Do("DEL", key); err != nil {
		return errors.Trace(err)
	}
	return nil
}

func KeyExists(conn redis.Conn, key string) (bool, error) {
	ok, err := redis.Bool(conn.Do("EXISTS", key))
	return ok, errors.Trace(err)
}

func SetStringList(conn redis.Conn, key string, list []string) error {
	conn.Send("DEL", key)
	l := len(list)
	if l > 0 {
		args := make([]interface{}, l+1)
		args[0] = key
		for i, s := range list {
			args[i+1] = s
		}
		conn.Send("RPUSH", args...)
	}
	conn.Flush()
	for i := 0; i < 2; i++ {
		if _, err := conn.Receive(); err != nil {
			return errors.Trace(err)
		}
	}
	return nil
}

func AppendStringList(conn redis.Conn, key string, list []string) error {
	l := len(list)
	if l > 0 {
		args := make([]interface{}, l+1)
		args[0] = key
		for i, s := range list {
			args[i+1] = s
		}
		conn.Send("RPUSH", args...)
	}
	conn.Flush()
	if _, err := conn.Receive(); err != nil {
		return errors.Trace(err)
	}
	return nil
}

func GetStringList(conn redis.Conn, key string) ([]string, error) {
	return redis.Strings(conn.Do("LRANGE", key, 0, -1))
}

func SetString(conn redis.Conn, key string, value string) error {
	err := conn.Send("SET", key, value)
	if err != nil {
		return errors.Trace(err)
	}
	err = conn.Flush()
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

func SetStringEx(conn redis.Conn, key string, value string, expireSeconds int) error {
	err := conn.Send("SETEX", key, expireSeconds, value)
	if err != nil {
		return errors.Trace(err)
	}
	err = conn.Flush()
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

func GetString(conn redis.Conn, key string) (string, error) {
	s, err := redis.String(conn.Do("GET", key))
	if err != nil {
		if err == redis.ErrNil {
			return "", nil
		}
		return "", errors.Trace(err)
	}
	return s, nil
}

func SetIntEx(conn redis.Conn, key string, value int, expireSeconds int) error {
	err := conn.Send("SETEX", key, expireSeconds, value)
	if err != nil {
		return errors.Trace(err)
	}
	err = conn.Flush()
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

func SetInt(conn redis.Conn, key string, value int) error {
	err := conn.Send("SET", key, value)
	if err != nil {
		return errors.Trace(err)
	}
	err = conn.Flush()
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

func Decr(conn redis.Conn, key string) error {
	err := conn.Send("DECR", key)
	if err != nil {
		return errors.Trace(err)
	}
	err = conn.Flush()
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

func GetInt(conn redis.Conn, key string) (int, error) {
	s, err := redis.Int(conn.Do("GET", key))
	if err != nil {
		if err == redis.ErrNil {
			return 0, nil
		}
		return 0, errors.Trace(err)
	}
	return s, nil
}

func TTL(conn redis.Conn, key string) (int64, error) {
	s, err := redis.Int64(conn.Do("TTL", key))
	if err != nil {
		if err == redis.ErrNil {
			return 0, nil
		}
		return 0, errors.Trace(err)
	}
	return s, nil
}

func GetStringWithoutNil(conn redis.Conn, key string) (string, bool, error) {
	result, err := redis.String(conn.Do("GET", key))
	if err == redis.ErrNil {
		return "", false, nil
	}
	if err != nil {
		return "", false, errors.Trace(err)
	}
	return result, true, nil
}

func Get(c redis.Conn, key string) ([]byte, error) {
	v, err := redis.Bytes(c.Do("GET", key))
	if err == redis.ErrNil {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Trace(err)
	}
	return v, nil
}

func Set(conn redis.Conn, key string, value []byte, expireTime int64) error {
	conn.Send("SET", key, value)
	conn.Send("EXPIRE", key, expireTime)
	conn.Flush()
	_, err := conn.Receive()
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

func MSet(conn redis.Conn, keys, values []interface{}, expireTime int64) error {
	luaScript := `for index,key in ipairs(KEYS) do redis.call('SETEX', key, ARGV[1], ARGV[index+1]) end`
	script := redis.NewScript(len(keys), luaScript)
	args := make([]interface{}, 0, len(keys)+len(values)+1)
	args = append(args, keys...)
	args = append(args, expireTime)
	args = append(args, values...)
	_, err := script.Do(conn, args...)
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

func MGet(c redis.Conn, keys []string) ([][]byte, error) {
	if len(keys) == 0 {
		return nil, nil
	}
	args := make([]interface{}, 0, len(keys))
	for _, f := range keys {
		args = append(args, f)
	}
	vals, err := redis.ByteSlices(c.Do("MGET", args...))
	if err != nil {
		return nil, errors.Trace(err)
	}
	return vals, nil
}

func SetEx(conn redis.Conn, keys []string, vals [][]byte, expire int64) error {
	if len(keys) == 0 || len(keys) != len(vals) {
		return fmt.Errorf("wrong number of arguments for SETEX")
	}
	for i, k := range keys {
		err := conn.Send("SETEX", k, expire, vals[i])
		if err != nil {
			return errors.Trace(err)
		}
	}
	err := conn.Flush()
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

func Del(conn redis.Conn, keys []string) error {
	if len(keys) == 0 {
		return nil
	}
	args := make([]interface{}, 0, len(keys))
	for _, k := range keys {
		args = append(args, k)
	}
	err := conn.Send("DEL", args...)
	if err != nil {
		return errors.Trace(err)
	}
	err = conn.Flush()
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

func GetHashKeys(conn redis.Conn, key string) ([]interface{}, error) {
	values, err := redis.Values(conn.Do("HKEYS", key))
	if err != nil {
		return nil, errors.Trace(err)
	}
	return values, nil
}

func SetNx(conn redis.Conn, key string, val string, expire int64) error {
	res, err := redis.Int(conn.Do("SETNX", key, val)) // res 1：设置成功 0： 设置失败，旧的值存在
	if err != nil {
		return err
	}
	if res == 1 {
		// 如果设置新的值成功了, 设置一个过期时间
		err := conn.Send("EXPIRE", key, expire)
		if err != nil {
			return err
		}
	}
	return nil
}
