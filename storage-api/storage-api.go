package storageapi

import (
	"encoding/binary"
	"encoding/json"
	"hello/bucket"
	"hello/iolayer"
	"hello/utils"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)


func RegisterStorageApiEndpoints(){
	http.Handle("/file", utils.PostMethod(writeFileToBucket))
}


func writeFileToBucket(w http.ResponseWriter, r *http.Request){
	multipartHeader := r.Header.Get("Content-Type") 

	if  multipartHeader == "" || !strings.Contains(multipartHeader, "multipart/form-data") {
		utils.WriteBasicResp(w, nil, 2, "Header multipart/form-data must be supplied")
		return
	}

	fileReader, fileHeader, err:= r.FormFile("file")

	if err != nil {
		utils.WriteBasicResp(w, nil, 3, "Form file is absent")
		return
	}



	fileName := fileHeader.Filename
	fileBytes, err := io.ReadAll(fileReader)
	
	if err != nil {
		utils.WriteBasicResp(w, nil, 4, "Failed extracting file content")
		return
	}

	if strings.Contains(fileName ,"/"){
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
			RawFile: &fileBytes,
			Filename: fileName,
		},
	}

	fileRef, err := ioLayer.SaveFile()

	if err != nil {
		utils.WriteBasicResp(w, nil, 1, err.Error())
		return
	}

	utils.WriteBasicResp(w, fileRef, 0, "")
}



func readFileFromStorage (w http.ResponseWriter, r *http.Request){
	bytes, err := os.ReadFile("/Users/markstepanov/go_stuff/hello/static/fistBucket/2025-07-14 08-09-38.mov/data.xxl")
	if err != nil {
		return
	}

	a := int(binary.BigEndian.Uint32(bytes[3:7]))
	jsonBytes := bytes[7: 7+a]

	myMap := map[string]interface{}{}
	json.Unmarshal(jsonBytes, &myMap)
	log.Println(myMap)

}

