package service

import (
	"github.com/jmoiron/sqlx/types"
	"github.com/nmmugia/marvel-character/models"
)

type (
	CharacterUsecase interface {
		GetCharacters(param models.GetCharactersParam) (result types.JSONText, errx models.Errorx)
		GetCharacterByIDs(ids string) (result types.JSONText, errx models.Errorx)
		GetAllDataByCache() (errx models.Errorx)
	}
)
