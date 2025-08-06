package server

import (
	"gofiel/bucket"
	"gofiel/config"
	storageapi "gofiel/storage-api"
	"log"
	"net/http"
)

func ServerStart() {
	err := config.ReadConfigFile()

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Config file read and set")

	err = bucket.RegisterBuckets()

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Bucket info read and set")

	registerEndpoints()

	log.Println("Server started on port: " + config.GlobalServerConfig.Port)

	log.Fatal(http.ListenAndServe(":"+config.GlobalServerConfig.Port, nil))
}

func registerEndpoints() {
	storageapi.RegisterStorageApiEndpoints()
	bucket.RegisterBucketApiEndpoints()
}
