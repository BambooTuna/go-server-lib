package redis

import (
	"fmt"
	"github.com/BambooTuna/quest-market/backend/settings"
	"github.com/go-redis/redis"
)

func RedisConnection(db int) *redis.Client {
	redisAddr := fmt.Sprintf("%s:%s",
		settings.FetchEnvValue("REDIS_HOST", "127.0.0.1"),
		settings.FetchEnvValue("REDIS_PORT", "6379"),
	)
	connection := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: settings.FetchEnvValue("REDIS_PASSWORD", ""),
		DB:       db,
	})
	return connection
}
