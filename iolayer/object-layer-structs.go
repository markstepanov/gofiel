package iolayer

import (
	"hello/bucket"
	"hello/comprassion"
	"os"
)

type ObjectFile struct {
	RawFile *[]byte
	CompressedFile *[]byte
	Filename string
	File * os.File
	ContentType string
	ComprassionInfo comprassion.ComprassionInfo
}

type IoLayer struct {
	Bucket *bucket.Bucket
	ObjectFile ObjectFile
}



type FileRef struct {
	Comprassion comprassion.ComprassionInfo
	Bucket string
}


type IoLayerApi interface{
	InitializeObjectStorage() (string, error)
	CompressFile() error
	SaveFile() (FileRef, error)
}


