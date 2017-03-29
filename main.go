package main

import (
	"log"
	"net/http"
	"os"

	"github.com/kashmire/greenHR/routing"
	"github.com/kashmire/greenHR/storage"
)

func getPort() string {
	port := os.Getenv("APP_PORT")
	if port == "" {
		log.Fatal("POST env is missing.")
	}

	return port
}

func main() {
	db := storage.Connection()
	mux := routing.Router(db)

	http.ListenAndServe(getPort(), mux)
}
