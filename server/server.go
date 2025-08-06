package server

import (
	"gofiel/config"
	storageapi "gofiel/storage-api"
	"log"
	"net/http"
)

func ServerStart() {

	config.ReadConfigFile()
	registerEndpoints()

	// TODO read configs

	log.Println("Server started on port 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func registerEndpoints() {
	storageapi.RegisterStorageApiEndpoints()
}
