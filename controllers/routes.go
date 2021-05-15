package controllers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nmmugia/marvel/service"
)

// ListRoutes is the function where you could set up url
func ListRoutes(router *mux.Router, characterUsecase service.CharacterUsecase) {
	handler := ServiceHandler(characterUsecase)
	router.HandleFunc("/characters", handler.GetCharacters).Methods("GET")
	router.HandleFunc("/characters/{ids}", handler.GetCharacterByIDs).Methods("GET")
	router.HandleFunc("/jobs/get-data-by-cache", handler.GetAllDataByCacheHourly).Methods("POST")
	router.PathPrefix("/docs/").Handler(http.StripPrefix("/docs/", http.FileServer(http.Dir("./documentation"))))
}
