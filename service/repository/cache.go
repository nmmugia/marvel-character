package repository

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/nmmugia/marvel/service"
)

type cacheRepository struct {
	redisClient *redis.Client
}

func NewCacheRepository(redisClient *redis.Client) service.CacheRepository {
	return cacheRepository{redisClient}
}

func (repo cacheRepository) Set(key string, value interface{}, duration time.Duration) error {
	return repo.redisClient.Set(key, value, duration).Err()
}

func (repo cacheRepository) Get(key string) (data string) {
	data, _ = repo.redisClient.Get(key).Result()
	return data
}

func (repo cacheRepository) GetAllKeysByPattern(pattern string) (data []string) {
	return repo.redisClient.Keys(pattern).Val()
}
