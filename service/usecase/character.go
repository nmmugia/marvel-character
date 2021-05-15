package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jmoiron/sqlx/types"
	"github.com/nmmugia/marvel-character/models"
	"github.com/nmmugia/marvel-character/service"
	u "github.com/nmmugia/marvel-character/utils"
)

type characterUsecase struct {
	cacheRepository service.CacheRepository
}

func NewCharacterUsecase(cacheRepository service.CacheRepository) service.CharacterUsecase {
	return characterUsecase{cacheRepository: cacheRepository}
}

func (uc characterUsecase) GetCharacters(param models.GetCharactersParam) (result types.JSONText, errx models.Errorx) {
	var paramURL string
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
	var cacheResult = uc.cacheRepository.Get("characters/" + paramURL)

	// guardian if cache exist
	if len(cacheResult) > 0 {
		err := json.Unmarshal([]byte(cacheResult), &result)
		// if no error found then return the data because the cache data is valid
		if err == nil {
			return result, errx
		}

		// if err not nil then jump to default flow
		log.Println("Error while marshalling result from cache, cause:" + err.Error())
	}

	// default flow

	resp, err := uc.HitMarvelsEndpoint("GET", "characters", paramURL)
	if err != nil {
		errx = models.CreateErrx(http.StatusInternalServerError, err, "While use HitMarvelsEndpoint")
		return result, errx
	}
	if len(resp.Message) > 0 {
		errx = models.CreateErrx(http.StatusInternalServerError, errors.New(resp.Message), "While use HitMarvelsEndpoint")
		return result, errx
	}

	errx.Status = resp.Code
	err = uc.cacheRepository.Set("characters/"+paramURL, string(resp.Data), time.Hour)
	if err != nil {
		log.Println("Error on GetCharacters func, while set cache on characters data. Retrying...")
		go func() {
			iteration := 1
			for {
				log.Println("set cache on characters data. Retrying attempt " + fmt.Sprint(iteration))
				if iteration >= 4 {
					log.Println("Failed to set cache on characters data")
					break
				}
				err = uc.cacheRepository.Set("characters/"+paramURL, string(resp.Data), time.Hour)
				if err != nil {
					time.Sleep(1 * time.Second)
					iteration++
					continue
				}
				log.Println("Success set cache on characters data attempt" + fmt.Sprint(iteration))
				break
			}
		}()
	}

	return resp.Data, errx
}

func (uc characterUsecase) GetCharacterByIDs(ids string) (result types.JSONText, errx models.Errorx) {
	cacheResult := uc.cacheRepository.Get("characters/" + ids)

	// guardian if cache exist
	if len(cacheResult) > 0 {
		err := json.Unmarshal([]byte(cacheResult), &result)
		// if no error found then return the data because the cache data is valid
		if err == nil {
			return result, errx
		}

		// if err not nil then jump to default flow
		log.Println("Error while marshalling result from cache, cause:" + err.Error())
	}

	resp, err := uc.HitMarvelsEndpoint("GET", "characters/"+ids, "")
	if err != nil {
		errx = models.CreateErrx(http.StatusInternalServerError, err, "While use HitMarvelsEndpoint")
		return result, errx
	}
	if len(resp.Message) > 0 {
		errx = models.CreateErrx(http.StatusInternalServerError, errors.New(resp.Message), "While use HitMarvelsEndpoint")
		return result, errx
	}
	errx.Status = resp.Code
	err = uc.cacheRepository.Set("characters/"+ids, string(resp.Data), time.Hour)
	if err != nil {
		log.Println("Error on GetCharacterByIDs func, while set cache on characters data. Retrying...")
		go func() {
			iteration := 1
			for {
				log.Println("set cache on characters data. Retrying attempt " + fmt.Sprint(iteration))
				if iteration >= 4 {
					log.Println("Failed to set cache on characters data")
					break
				}
				err = uc.cacheRepository.Set("characters/"+ids, string(resp.Data), time.Hour)
				if err != nil {
					time.Sleep(1 * time.Second)

					iteration++
					continue
				}
				log.Println("Success set cache on characters data attempt" + fmt.Sprint(iteration))
				break
			}
		}()
	}
	return resp.Data, errx
}

func (uc characterUsecase) GetAllDataByCache() (errx models.Errorx) {
	cacheResult := uc.cacheRepository.GetAllKeysByPattern("characters*")

	// guardian if cache exist
	if len(cacheResult) > 0 {
		for _, v := range cacheResult {
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

			resp, err := uc.HitMarvelsEndpoint("GET", url, param)
			if err != nil {
				errx = models.CreateErrx(http.StatusInternalServerError, err, "While use HitMarvelsEndpoint")
				log.Println(errx.Message, errx.Err.Error())
			}
			if len(resp.Message) > 0 {
				errx = models.CreateErrx(http.StatusInternalServerError, errors.New(resp.Message), "While use HitMarvelsEndpoint")
				log.Println(errx.Message, errx.Err.Error())
			}
			errx.Status = resp.Code
			err = uc.cacheRepository.Set(v, string(resp.Data), time.Hour)
			if err != nil {
				log.Println("Error on GetCharacterByIDs func, while set cache on characters data. Retrying...")
				go func() {
					iteration := 1
					for {
						log.Println("set cache on characters data. Retrying attempt " + fmt.Sprint(iteration))
						if iteration >= 4 {
							log.Println("Failed to set cache on characters data")
							break
						}
						err = uc.cacheRepository.Set(v, string(resp.Data), time.Hour)
						if err != nil {
							time.Sleep(1 * time.Second)

							iteration++
							continue
						}
						log.Println("Success set cache on characters data attempt" + fmt.Sprint(iteration))
						break
					}
				}()
			}
		}
	}
	errx.Status = http.StatusOK
	errx.Message = "Cron has been successfully executed"
	return errx
}

func (uc characterUsecase) HitMarvelsEndpoint(method string, path string, params string) (result models.MarvelsResult, err error) {
	var (
		ts  = time.Now().Unix()
		url = fmt.Sprintf("%s/%s?ts=%d&apikey=%s&hash=%s",
			strings.TrimRight(os.Getenv("MARVEL_BASE_URL"), "/"),
			strings.TrimLeft(path, "/"), ts,
			os.Getenv("PUBLIC_KEY"),
			u.StringToMD5(fmt.Sprint(ts)+os.Getenv("PRIVATE_KEY")+os.Getenv("PUBLIC_KEY")),
		)
		client = &http.Client{}
	)
	req, err := http.NewRequest(method, url+params, nil)
	if err != nil {
		return result, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	if err := json.Unmarshal([]byte(bodyBytes), &result); err != nil {
		return result, err
	}
	return result, err
}
