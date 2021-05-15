package config

import (
	"log"
	"os"

	"github.com/go-redis/redis"
	"github.com/nmmugia/marvel-character/service/usecase"

	"github.com/nmmugia/marvel-character/controllers"

	"github.com/gorilla/mux"
	"github.com/nmmugia/marvel-character/service/repository"
)

type Config struct {
	Redis *redis.Client
	Route *mux.Router
}

func (cfg *Config) InitService() {
	router := mux.NewRouter()
	cacheRepo := repository.NewCacheRepository(cfg.Redis)
	characterUsecase := usecase.NewCharacterUsecase(cacheRepo)

	controllers.ListRoutes(router, characterUsecase)
	controllers.ListJobs(characterUsecase)
	cfg.Route = router

}

func (cfg *Config) InitCache() {
	cfg.Redis = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_CACHE_URL"),
		Password: os.Getenv("REDIS_CACHE_PWD"),
		DB:       0,
	})

	strStatus, err := cfg.Redis.Ping().Result()
	log.Println(strStatus, err)
}
