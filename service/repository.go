package service

import (
	"time"
)

type (
	CacheRepository interface {
		Get(key string) (result string)
		Set(key string, value interface{}, duration time.Duration) (err error)
		GetAllKeysByPattern(pattern string) (data []string)
	}
)
