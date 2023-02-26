package redis

import (
	"github.com/gomodule/redigo/redis"
	"github.com/obgnail/plugin-platform/common/config"
	"github.com/obgnail/plugin-platform/common/errors"
	"time"
)

const (
	defaultRedisPost        = "localhost:6379"
	defaultRedisPassword    = ""
	defaultRedisMaxIdle     = 3
	defaultRedisMaxActive   = 20
	defaultRedisIdleTimeout = 180
	defaultRedisExpire      = 2592000 // 1*30*24*60*60
)

var (
	pluginRedisPool *redis.Pool
	redisKeyExpire  int
)

func InitRedis() error {
	addr := config.String("platform.redis_address", defaultRedisPost)
	password := config.String("platform.redis_password", defaultRedisPassword)
	maxIdle := config.Int("platform.redis_max_idle", defaultRedisMaxIdle)
	maxActive := config.Int("platform.redis_max_active", defaultRedisMaxActive)
	idleTimeout := config.Int("platform.redis_idle_timeout", defaultRedisIdleTimeout)
	redisKeyExpire = config.Int("platform.redis_expire", defaultRedisExpire)
	pluginDB := config.IntOrPanic("platform.redis_plugin_db")

	initRedisPoolFunc := func(db int) *redis.Pool {
		return &redis.Pool{
			MaxIdle:     maxIdle,
			MaxActive:   maxActive,
			IdleTimeout: time.Duration(idleTimeout) * time.Second,
			Dial: func() (redis.Conn, error) {
				var c redis.Conn
				var err error
				if password == "" {
					c, err = redis.Dial("tcp", addr)
				} else {
					c, err = redis.Dial("tcp", addr, redis.DialPassword(password))
				}
				if err != nil {
					return nil, errors.Trace(err)
				}
				c.Do("SELECT", db)
				return c, nil
			},
		}
	}
	pluginRedisPool = initRedisPoolFunc(pluginDB)
	return nil
}

func GetRedisConn() redis.Conn { return pluginRedisPool.Get() }
