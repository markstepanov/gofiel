package server

import (
	storageapi "hello/storage-api"
	"log"
	"net/http"
)

func ServerStart(){

	registerEndpoints()

	// TODO read configs

	log.Println("Server started on port 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}


func registerEndpoints(){
	storageapi.RegisterStorageApiEndpoints()
}