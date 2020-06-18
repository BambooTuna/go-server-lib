package redis

import (
	"fmt"
	"github.com/BambooTuna/go-server-lib/config"
	goredis "github.com/go-redis/redis"
)

func RedisConnection(db int) *goredis.Client {
	redisAddr := fmt.Sprintf("%s:%s",
		config.GetEnvString("REDIS_HOST", "127.0.0.1"),
		config.GetEnvString("REDIS_PORT", "6379"),
	)
	connection := goredis.NewClient(&goredis.Options{
		Addr:     redisAddr,
		Password: config.GetEnvString("REDIS_PASSWORD", ""),
		DB:       db,
	})
	return connection
}
