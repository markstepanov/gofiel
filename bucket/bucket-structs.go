package bucket

import (
	"gofiel/comprassion"
)

type Bucket struct {
	Id              int
	Name            string
	Path            string
	CompressionType string
}

type BucketResp struct {
	Id              int
	Name            string
	CompressionType string
}

func (bucket *Bucket) toBucketResp() BucketResp {
	return BucketResp{bucket.Id, bucket.Name, bucket.CompressionType}
}

var testBucket1 = Bucket{
	Id:              1,
	Name:            "fistBucket",
	Path:            "/Users/markstepanov/go_stuff/hello/static/fistBucket",
	CompressionType: comprassion.ComressionZstd,
}

var tetsbucket2 = Bucket{
	Id:              2,
	Name:            "secondBucket",
	Path:            "/Users/markstepanov/go_stuff/hello/static/secondBucket",
	CompressionType: comprassion.ComressionGzip,
}

var Bukets = []Bucket{testBucket1, tetsbucket2}
