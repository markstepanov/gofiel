package iolayer

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"gofiel/comprassion"
	"os"
	"path"
	"strings"
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
		UncompressedSize:     len(*ioLayer.ObjectFile.RawFile),
		CompressedSize:       len(*ioLayer.ObjectFile.CompressedFile),
		ComprassionAlgorithm: ioLayer.Bucket.CompressionType,
	}

	err = writeToFile(ioLayer)

	if err != nil {
		// TODO:  decide on logic that will delete whole object storage files and path
		return nil, err

	}

	ioLayer.ObjectFile.File.Close()

	return &FileRef{
		Comprassion: ioLayer.ObjectFile.ComprassionInfo,
		Bucket:      ioLayer.Bucket.Name,
	}, nil
}

func (ioLayer *IoLayer) FindFile() error {
	objectPath := path.Join(ioLayer.Bucket.Path, ioLayer.ObjectFile.Filename)

	err := verifyGetFileRequest(ioLayer, objectPath)

	if err != nil {
		return err
	}

	filePath := path.Join(objectPath, "data.xxl")
	bytes, err := os.ReadFile(path.Join(objectPath, "data.xxl"))

	if err != nil {
		return errors.New("failed while reading filepath: " + filePath)
	}

	header := string(bytes[0:3])

	if header != "XXL" {
		return errors.New("file does not contain XXL header.filepath: " + filePath)
	}

	metaInfoLen := int(binary.BigEndian.Uint32(bytes[3:7]))
	jsonBytes := bytes[7 : 7+metaInfoLen]
	metaInfo := ObjectMetadata{}

	err = json.Unmarshal(jsonBytes, &metaInfo)

	if err != nil {
		return errors.New("failed to read metaInfo for file : " + filePath)
	}

	compressedFileObject := bytes[7+metaInfoLen:]

	if len(compressedFileObject) == 0 {
		return errors.New("header and metainfo is present, but compressedFileObject len is 0 for file: " + filePath)
	}

	decompressedBytes, err := comprassion.DecompresBytes(&compressedFileObject, &metaInfo.ComprassionInfo)

	if err != nil {
		return err
	}

	ioLayer.ObjectFile.CompressedFile = &compressedFileObject
	ioLayer.ObjectFile.ComprassionInfo = metaInfo.ComprassionInfo
	ioLayer.ObjectFile.RawFile = decompressedBytes
	ioLayer.ObjectFile.ContentType = metaInfo.ContentType

	return nil
}

func verifyGetFileRequest(ioLayer *IoLayer, objectPath string) error {
	if strings.Contains(ioLayer.ObjectFile.Filename, "/") {
		return errors.New("invalid filename")
	}

	pathInfo, err := os.Stat(objectPath)
	if err != nil {
		return err
	}

	if !pathInfo.IsDir() {
		return errors.New("unknown path")
	}
	return nil
}

func writeToFile(ioLayer *IoLayer) error {

	attrs := ObjectMetadata{
		ComprassionInfo: ioLayer.ObjectFile.ComprassionInfo,
		ContentType:     ioLayer.ObjectFile.ContentType,
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
