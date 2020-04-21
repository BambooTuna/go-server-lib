package session_v2

import (
	"github.com/go-redis/redis"
	"time"
)

type RedisSessionStorageDao struct {
	Client *redis.Client
}

func (r RedisSessionStorageDao) Store(key, value string, expiration time.Duration) (string, error) {
	return r.Client.Set(key, value, expiration).Result()
}

func (r RedisSessionStorageDao) Find(key string) (string, error) {
	result, err := r.Client.Get(key).Result()
	return result, err
}

func (r RedisSessionStorageDao) Remove(key string) (int64, error) {
	return r.Client.Del(key).Result()
}

func (r RedisSessionStorageDao) Refresh(key string, expiration time.Duration) (bool, error) {
	return r.Client.Expire(key, expiration).Result()
}
