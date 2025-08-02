package bucket

import "errors"

func FindBucketById(bucketId int) (*Bucket, error){
	for _, bucket := range Bukets {
		if bucket.Id == bucketId {
			return &bucket, nil
		}
	}

	return  nil, errors.New("could not find bucket by id")
}