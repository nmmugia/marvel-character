package repository

import (
	"log"
	"time"

	"github.com/go-redis/redis"
	"github.com/nmmugia/marvel-character/service"
)

type cacheRepository struct {
	redisClient redis.Cmdable
}

func NewCacheRepository(redisClient redis.Cmdable) service.CacheRepository {
	return cacheRepository{redisClient}
}

func (repo cacheRepository) Set(key string, value interface{}, duration time.Duration) error {
	if err := repo.redisClient.Set(key, value, duration).Err(); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (repo cacheRepository) Get(key string) (data string) {
	data, _ = repo.redisClient.Get(key).Result()
	return data
}

func (repo cacheRepository) GetAllKeysByPattern(pattern string) (data []string) {
	return repo.redisClient.Keys(pattern).Val()
}
