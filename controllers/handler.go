package controllers

import (
	"log"
	"net/http"

	"github.com/nmmugia/marvel-character/models"

	"github.com/gorilla/mux"

	"github.com/nmmugia/marvel-character/service"
	u "github.com/nmmugia/marvel-character/utils"
)

type Handler struct {
	characterUsecase service.CharacterUsecase
}

func ServiceHandler(characterUsecase service.CharacterUsecase) *Handler {
	return &Handler{
		characterUsecase: characterUsecase,
	}
}

func (handler Handler) GetCharacters(w http.ResponseWriter, r *http.Request) {
	var param = models.GetCharactersParam{
		Name:           r.URL.Query().Get("name"),
		NameStartsWith: r.URL.Query().Get("name_starts_with"),
		OrderBy:        r.URL.Query().Get("order_by"),
		Comics:         r.URL.Query().Get("comics"),
		Events:         r.URL.Query().Get("events"),
		Stories:        r.URL.Query().Get("stories"),
		Series:         r.URL.Query().Get("series"),
		Limit:          u.StringToInt(r.URL.Query().Get("limit"), 100),
		Offset:         u.StringToInt(r.URL.Query().Get("offset"), 0),
	}
	param.ModifiedSince, _ = u.ParseFromString(r.URL.Query().Get("modifiedSince"))
	resp, errx := handler.characterUsecase.GetCharacters(param)
	if errx.Err != nil {
		u.Response(w, nil, errx)
		return
	}
	u.Response(w, resp, errx)
}

func (handler Handler) GetCharacterByIDs(w http.ResponseWriter, r *http.Request) {
	var (
		vars = mux.Vars(r)
		ids  string
	)
	if v, ok := vars["ids"]; ok {
		ids = v
	}
	resp, errx := handler.characterUsecase.GetCharacterByIDs(ids)
	if errx.Err != nil {
		u.Response(w, nil, errx)
		return
	}
	u.Response(w, resp, errx)
}

func (handler Handler) GetAllDataByCacheHourly(w http.ResponseWriter, r *http.Request) {
	errx := handler.characterUsecase.GetAllDataByCache()
	log.Println(errx.Message)
	u.Response(w, nil, models.Errorx{})
}
