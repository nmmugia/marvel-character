package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/nmmugia/marvel/config"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	port := os.Getenv("PORT")
	if len(port) > 0 {
		port = "8080"
	}

	log.Println("Listening to port:", port)

	var server config.Config

	server.InitCache()
	server.InitService()

	err = http.ListenAndServe(":"+port, server.Route)
	if err != nil {
		panic(err)
	}
}
