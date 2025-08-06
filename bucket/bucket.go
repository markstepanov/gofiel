package bucket

import (
	"encoding/json"
	"errors"
	"gofiel/config"
	"log"
	"os"
	"path"
)

func FindBucketById(bucketId int) (*Bucket, error) {
	bucket, ok := Bukets[bucketId]
	if !ok {
		return nil, errors.New("could not find bucket by id")
	}
	return &bucket, nil
}

func RegisterBuckets() error {
	err := checkIfBasePathExists()

	if err != nil {
		return err
	}

	bucketList, err := readOrCreateBucketsMetaInfo()

	if err != nil {
		return err
	}

	if len(*bucketList) == 0 {
		log.Println("buckets.metainfo is not presnet! use POST /bucket to create your first bucket!")
	}

	return nil
}

func checkIfBasePathExists() error {
	pathInfo, err := os.Stat(config.GlobalServerConfig.Basedir)
	if err != nil {
		return err
	}

	if !pathInfo.IsDir() {
		return errors.New("base-path is not a directory")
	}

	return nil
}

func readOrCreateBucketsMetaInfo() (*[]Bucket, error) {

	bucketsMetaInfoPath := path.Join(config.GlobalServerConfig.Basedir, "buckets.metainfo")
	data, err := os.ReadFile(bucketsMetaInfoPath)

	if err != nil {
		err = createBucketsMetadataFile(bucketsMetaInfoPath)
		if err != nil {
			return nil, err
		}
	}

	bucketsMetaInfo := []Bucket{}
	if len(data) == 0 {
		return &bucketsMetaInfo, nil
	}

	err = json.Unmarshal(data, &bucketsMetaInfo)

	if err != nil {
		return nil, err
	}

	return &bucketsMetaInfo, nil
}

func createBucketsMetadataFile(bucketsMetaInfoPath string) error {
	_, err := os.Create(bucketsMetaInfoPath)
	if err != nil {
		return err
	}
	return nil
}
