package iolayer

import (
	"gofiel/bucket"
	"gofiel/comprassion"
	"os"
)

type ObjectMetadata struct {
	ComprassionInfo comprassion.ComprassionInfo
	ContentType     string
}

type ObjectFile struct {
	RawFile         *[]byte
	CompressedFile  *[]byte
	Filename        string
	File            *os.File
	ContentType     string
	ComprassionInfo comprassion.ComprassionInfo
}

type IoLayer struct {
	Bucket     *bucket.Bucket
	ObjectFile ObjectFile
}

type FileRef struct {
	Comprassion comprassion.ComprassionInfo
	Bucket      string
}

type IoLayerApi interface {
	InitializeObjectStorage() (string, error)
	CompressFile() error
	DecompressFIle() error
	SaveFile() (FileRef, error)
	FindFile() error
}
