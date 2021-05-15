package repository_test

import (
	"testing"
	"time"

	"github.com/nmmugia/marvel-character/service/repository"

	"github.com/elliotchance/redismock/v7"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
)

var (
	client *redis.Client
	key    = "key"
	val    = "val"
)

func TestSet(t *testing.T) {
	exp := time.Duration(0)

	mock := redismock.NewNiceMock(client)
	mock.On("Set", key, val, exp).Return(redis.NewStatusResult("", nil))

	r := repository.NewCacheRepository(mock)
	err := r.Set(key, val, exp)
	assert.NoError(t, err)
}

func TestGet(t *testing.T) {
	mock := redismock.NewNiceMock(client)
	mock.On("Get", key).Return(redis.NewStringResult(val, nil))

	r := repository.NewCacheRepository(mock)
	res := r.Get(key)
	assert.Equal(t, val, res)
}
