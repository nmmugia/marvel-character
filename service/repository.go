package service

import (
	"time"

	"github.com/nmmugia/marvel/models"
)

type (
	CharacterRepository interface {
		GetCharacters(param models.GetCharactersParam) (result map[string]interface{}, errx models.Errorx)
		GetCharacterByIDs(ids string) (result map[string]interface{}, errx models.Errorx)
	}

	CacheRepository interface {
		Get(key string) (result string)
		Set(key string, value interface{}, duration time.Duration) (err error)
		GetAllKeysByPattern(pattern string) (data []string)
	}
)
