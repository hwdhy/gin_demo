package base

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
)

// 初始化redis连接池
func initRedisPool() {
	config := GConf.RedisConfig
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	RedisPool = &redis.Pool{
		MaxIdle:     config.MaxIdle,
		MaxActive:   config.MaxActive,
		IdleTimeout: time.Duration(config.IdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(
				"tcp",
				addr,
				redis.DialReadTimeout(time.Duration(GConf.CommonRedisTimeout.RedisDialReadTimeout)*time.Millisecond),
				redis.DialWriteTimeout(time.Duration(GConf.CommonRedisTimeout.RedisDialWriteTimeout)*time.Millisecond),
				redis.DialConnectTimeout(time.Duration(GConf.CommonRedisTimeout.RedisDialConnectTimeout)*time.Millisecond),
				redis.DialDatabase(config.Db),
				redis.DialPassword(config.Pass),
			)
			if err != nil {
				return nil, err
			}
			return c, err
		},
	}
	c := RedisPool.Get()
	defer c.Close()
	_, err := c.Do("PING")
	if err != nil {
		MultipleLog.Fatalf("Redis init pool failed, addr: %s, error: %s", addr, err.Error())
	}
	MultipleLog.Info("Redis init pool success")
}
