package bucket

import (
	"encoding/json"
	"errors"
	"gofiel/comprassion"
	"gofiel/config"
	"log"
	"os"
	"path"
)

func FindBucketByName(bucketId string) (*Bucket, error) {
	bucket, ok := Buckets[bucketId]
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


	for _, bucket := range *bucketList {
		Buckets[bucket.Name] = bucket
	}


	if err != nil {
		return err
	}


	if len(*bucketList) == 0 {
		log.Println("buckets.metainfo is not presnet! use POST /bucket to create your first bucket!")
	}

	return nil
}

// TODO: Validate existing buckets

// func validateActualBuckets() error{
// 	dirs, err := os.ReadDir(config.GlobalServerConfig.Basedir)
//
// 	if err != nil {
// 		return  err
// 	}
// 	for _, dir := range dirs {
// 		log.
// 	}
//
// }

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

func createNewBucket(bucketName string) (*Bucket, error) {
	err := checkIfBucketExists(bucketName)

	if err != nil {
		return nil, err
	}

	newBucketPath := path.Join(config.GlobalServerConfig.Basedir, bucketName)

	_, err = os.Stat(newBucketPath)

	if err == nil {
		return nil, errors.New("bucket " + bucketName + " is already exists")
	}

	err = os.Mkdir(newBucketPath, 0755)

	if err != nil {
		return nil, err
	}


	newBucket := Bucket{
		Name:            bucketName,
		Path:            newBucketPath,
		CompressionType: comprassion.ComressionZstd,
	}

	err = addBucketToBucketsMetaInfo(&newBucket)

	if err != nil {
		return nil, errors.New("failed creating bucket with name:" + bucketName)
	}

	addNewBucketToBucketCache(&newBucket)

	return &newBucket, nil
}

func addNewBucketToBucketCache(newBucket *Bucket) {
	Buckets[newBucket.Name] = *newBucket
}

func addBucketToBucketsMetaInfo(bucket *Bucket) error {
	bucketsMetaInfoPath := path.Join(config.GlobalServerConfig.Basedir, "buckets.metainfo")
	file, err := os.OpenFile(bucketsMetaInfoPath, os.O_RDWR, 0644)

	if err != nil {
		return err
	}

	defer file.Close()


	fileInfo , err :=  file.Stat()
	if err != nil {
		return errors.New("failed while geting file Stat")
	}

	existingBuckets := []Bucket{}

	if fileInfo.Size() != 0 {
		err = json.NewDecoder(file).Decode(&existingBuckets)
	}

	if err != nil {
		return err
	}

	existingBuckets = append(existingBuckets, *bucket)

	if err := file.Truncate(0); err != nil {
		return err
	}

	if _, err := file.Seek(0, 0); err != nil {
		return err
	}

	return json.NewEncoder(file).Encode(existingBuckets)
}

func checkIfBucketExists(newBucket string) error {
	_, err := FindBucketByName(newBucket)

	if err == nil {
		return errors.New("bucket with name" + newBucket + " is already exists")
	}

	return nil
}
