package storageapi

import (
	"gofiel/bucket"
	"gofiel/iolayer"
	"gofiel/utils"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
)

func RegisterStorageApiEndpoints() {
	http.HandleFunc("/file", handleFileEndoint)
}

func handleFileEndoint(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		writeFileToBucket(w, r)
	case http.MethodGet:
		getFileFromBucket(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func writeFileToBucket(w http.ResponseWriter, r *http.Request) {
	multipartHeader := r.Header.Get("Content-Type")

	if multipartHeader == "" || !strings.Contains(multipartHeader, "multipart/form-data") {
		utils.WriteBasicResp(w, nil, 2, "Header multipart/form-data must be supplied")
		return
	}

	fileReader, fileHeader, err := r.FormFile("file")

	if err != nil {
		utils.WriteBasicResp(w, nil, 3, "Form file is absent")
		return
	}

	fileName := fileHeader.Filename
	fileContentType := getContentTypeFromPart(fileHeader)
	fileBytes, err := io.ReadAll(fileReader)

	if err != nil {
		utils.WriteBasicResp(w, nil, 4, "Failed extracting file content")
		return
	}

	if strings.Contains(fileName, "/") {
		utils.WriteBasicResp(w, nil, 5, "Symbol '/' is not allowed in the filename")
		return
	}

	bucketid, err := strconv.Atoi(r.PostForm.Get("bucket-id"))

	if err != nil {
		utils.WriteBasicResp(w, nil, 6, "Bucket form is not present")
		return
	}

	bucket, err := bucket.FindBucketById(bucketid)

	if err != nil {
		utils.WriteBasicResp(w, nil, 4, "Bucket is not present")
		return
	}

	var ioLayer iolayer.IoLayer = iolayer.IoLayer{
		Bucket: bucket,
		ObjectFile: iolayer.ObjectFile{
			RawFile:     &fileBytes,
			Filename:    fileName,
			ContentType: fileContentType,
		},
	}

	fileRef, err := ioLayer.SaveFile()

	if err != nil {
		utils.WriteBasicResp(w, nil, 1, err.Error())
		return
	}

	utils.WriteBasicResp(w, fileRef, 0, "")
}

func getFileFromBucket(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("filename")

	if filename == "" {
		utils.WriteBasicResp(w, nil, 1, "filename is not specified")
		return
	}

	bucketid, err := strconv.Atoi(r.Header.Get("bucket-id"))

	if err != nil {
		utils.WriteBasicResp(w, nil, 6, "Bucket form is not present")
		return
	}

	bucket, err := bucket.FindBucketById(bucketid)

	if err != nil {
		utils.WriteBasicResp(w, nil, 4, "Bucket is not present")
		return
	}

	ioLayer := iolayer.IoLayer{
		Bucket: bucket,
		ObjectFile: iolayer.ObjectFile{
			Filename: filename,
		},
	}

	err = ioLayer.FindFile()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Write([]byte("Hello world"))
}

func getContentTypeFromPart(header *multipart.FileHeader) string {
	if ct := header.Header.Get("Content-Type"); ct != "" {
		return ct
	}
	return "application/octet-stream"
}
