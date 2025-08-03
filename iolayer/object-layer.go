package iolayer

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"gofiel/comprassion"
	"os"
	"path"
)

func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false // not found or inaccessible
	}
	return info.IsDir()
}

func (ioLayer *IoLayer) InitializeObjectStorage() error {
	objectPath := path.Join(ioLayer.Bucket.Path, ioLayer.ObjectFile.Filename)

	// TODO right now just replace spaces to undescore
	// objectPath = strings.ReplaceAll(objectPath, " ", "_")

	if dirExists(objectPath) {
		return errors.New("object already exists")
	}

	err := os.MkdirAll(objectPath, 0755)

	if err != nil {
		return err
	}

	objectPath = path.Join(objectPath, "data.xxl")

	file, err := os.Create(objectPath)

	if err != nil {
		return err
	}

	// file will be closed after all manipulations

	ioLayer.ObjectFile.File = file

	return nil
}

func (ioLayer *IoLayer) CompressFile() error {
	compressedFile, err := comprassion.CompressBytes(ioLayer.ObjectFile.RawFile, ioLayer.Bucket.CompressionType)
	if err != nil {
		return err
	}

	ioLayer.ObjectFile.CompressedFile = compressedFile
	return nil
}

func (ioLayer *IoLayer) DecompressFIle() error {
	return nil
}

func (ioLayer *IoLayer) SaveFile() (*FileRef, error) {
	err := ioLayer.InitializeObjectStorage()

	if err != nil {
		return nil, err
	}

	err = ioLayer.CompressFile()
	if err != nil {
		return nil, err
	}

	ioLayer.ObjectFile.ComprassionInfo = comprassion.ComprassionInfo{
		UncompressedSize: len(*ioLayer.ObjectFile.RawFile),
		CompressedSize:   len(*ioLayer.ObjectFile.CompressedFile),
	}

	err = writeToFile(ioLayer)

	if err != nil {
		// todo decide on logic that will delete whole object storage files and path
		return nil, err

	}

	ioLayer.ObjectFile.File.Close()

	return &FileRef{
		Comprassion: ioLayer.ObjectFile.ComprassionInfo,
		Bucket:      ioLayer.Bucket.Name,
	}, nil
}

func (ioLayer *IoLayer) FindFile() error {
	// TODO currenly working on this
	// bucketId := r.Header.Get("bucket-id")

	bytes, err := os.ReadFile("/Users/markstepanov/go_stuff/hello/static/fistBucket/2025-07-14 08-09-38.mov/data.xxl")
	if err != nil {
		return err
	}

	a := int(binary.BigEndian.Uint32(bytes[3:7]))
	jsonBytes := bytes[7 : 7+a]

	myMap := map[string]any{}
	json.Unmarshal(jsonBytes, &myMap)

	return nil
}

func writeToFile(ioLayer *IoLayer) error {

	attrs := map[string]any{
		"copressionAlgorithm": ioLayer.Bucket.CompressionType,
		"comprassionInfo":     ioLayer.ObjectFile.ComprassionInfo,
		"contentType":         ioLayer.ObjectFile.ContentType,
	}

	jsonBytes, err := json.Marshal(attrs)
	if err != nil {
		return err
	}

	length := uint32(len(jsonBytes))

	jsonLen := make([]byte, 4)
	binary.BigEndian.PutUint32(jsonLen, length)

	header := append(
		[]byte{},
		[]byte("XXL")...,
	)

	header = append(header, jsonLen...)
	header = append(header, jsonBytes...)
	ioLayer.ObjectFile.File.Write(header)
	ioLayer.ObjectFile.File.Write(*ioLayer.ObjectFile.CompressedFile)

	return nil
}
