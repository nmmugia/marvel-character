package usecase_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jmoiron/sqlx/types"
	"github.com/nmmugia/marvel-character/models"
	u "github.com/nmmugia/marvel-character/utils"

	"github.com/nmmugia/marvel-character/service/usecase"

	"github.com/bxcodec/faker"
	"github.com/nmmugia/marvel-character/service/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetCharacters(t *testing.T) {
	var (
		mockCache     = new(mocks.CacheRepository)
		mockCharacter = new(mocks.CharacterUsecase)
		r             = usecase.NewCharacterUsecase(mockCache)
		param         models.GetCharactersParam
		result        models.MarvelsResult
		resultCache   = `{"data": "spiderman"}`
		paramURL      string
	)

	err := faker.FakeData(&param)
	if err != nil {
		t.Fatalf("error while fake data %+v", err)
	}
	if len(param.Comics) > 0 {
		paramURL += fmt.Sprintf("&comics=%s", param.Comics)
	}
	if len(param.Events) > 0 {
		paramURL += fmt.Sprintf("&events=%s", param.Events)
	}
	if param.Limit > 0 {
		paramURL += fmt.Sprintf("&limit=%d", param.Limit)
	} else {
		paramURL += fmt.Sprintf("&limit=%d", 0)
	}
	if !param.ModifiedSince.IsZero() {
		paramURL += fmt.Sprintf("&modifiedSince=%s", param.ModifiedSince.Format("2006-01-02"))
	}
	if len(param.Name) > 0 {
		paramURL += fmt.Sprintf("&name=%s", param.Name)
	}
	if len(param.NameStartsWith) > 0 {
		paramURL += fmt.Sprintf("&nameStartsWith=%s", param.NameStartsWith)
	}
	if len(param.OrderBy) > 0 {
		paramURL += fmt.Sprintf("&orderBy=%s", param.OrderBy)
	}
	if len(param.Series) > 0 {
		paramURL += fmt.Sprintf("&series=%s", param.Series)
	}
	if len(param.Stories) > 0 {
		paramURL += fmt.Sprintf("&stories=%s", param.Stories)
	}
	paramURL += fmt.Sprintf("&offset=%d", param.Offset)
	mockCache.On("Get", "characters/"+paramURL).Return(resultCache, nil).Once()
	mockCharacter.On("HitMarvelsEndpoint", "GET", "characters", paramURL).Return(result, nil)
	mockCache.On("Set", "characters/"+paramURL).Return(nil).Once()

	res, errx := r.GetCharacters(param)
	assert.NoError(t, errx.Err)
	assert.NotEmpty(t, res)

}

func TestGetCharacterByIDs(t *testing.T) {
	var (
		mockCache     = new(mocks.CacheRepository)
		mockCharacter = new(mocks.CharacterUsecase)
		r             = usecase.NewCharacterUsecase(mockCache)
		resultCache   = `{"data": "spiderman"}`
		result        models.MarvelsResult
		ids           string
	)
	err := faker.FakeData(&ids)
	if err != nil {
		t.Fatalf("error while fake data %+v", err)
	}
	err = faker.FakeData(&result)
	if err != nil {
		t.Fatalf("error while fake data %+v", err)
	}
	mockCache.On("Get", "characters/"+ids).Return(resultCache, models.Errorx{Status: 200}).Once()
	mockCharacter.On("HitMarvelsEndpoint", "GET", "characters/"+ids, "").Return(result, nil)

	mockCache.On("Set", "characters/"+ids).Return(models.Errorx{Status: 200}).Once()

	res, errx := r.GetCharacterByIDs(ids)
	assert.NoError(t, errx.Err)
	assert.NotEmpty(t, res)
}

func TestGetAllDataByCache(t *testing.T) {
	var (
		mockCache     = new(mocks.CacheRepository)
		mockCharacter = new(mocks.CharacterUsecase)
		resultCache   = []string{`{"data": "spiderman"}`, `{"data": "bond, james bond"}`}
		result        models.MarvelsResult
	)
	err := faker.FakeData(&result)
	if err != nil {
		t.Fatalf("error while fake data %+v", err)
	}
	mockCache.On("GetAllKeysByPattern", "characters*").Return(resultCache, models.Errorx{Status: 200}).Once()
	if len(resultCache) > 0 {
		for _, v := range resultCache {
			var (
				vStrings = strings.Split(v, "/")
				url      = "characters"
				param    string
			)

			if len(vStrings) > 1 {
				url = "characters/" + vStrings[1]
			}

			if len(vStrings) > 1 && strings.Contains(vStrings[1], "&") {
				param = vStrings[1]
				url = "characters"
			}

			mockCharacter.On("HitMarvelsEndpoint", "GET", url, param).Return(result, nil).Once()

			mockCache.On("Set", v, string(types.JSONText{}), time.Hour).Return(models.Errorx{Status: 200}).Once()
		}
	}

}

func TestHitMarvelsEndpoint(t *testing.T) {
	var (
		ts  = time.Now().Unix()
		url = fmt.Sprintf("%s/%s?ts=%d&apikey=%s&hash=%s",
			strings.TrimRight(os.Getenv("MARVEL_BASE_URL"), "/"),
			strings.TrimLeft("character", "/"), ts,
			os.Getenv("PUBLIC_KEY"),
			u.StringToMD5(fmt.Sprint(ts)+os.Getenv("PRIVATE_KEY")+os.Getenv("PUBLIC_KEY")),
		)
		client = &http.Client{}
		result models.MarvelsResult
	)
	assert.Equal(t, time.Now().Unix(), ts)
	assert.IsType(t, "character/url", url)
	assert.IsType(t, &http.Client{Timeout: time.Minute}, client)
	err := json.Unmarshal([]byte(`{"data": "j. jonah jameson"}`), &result)
	assert.NoError(t, err)
}
