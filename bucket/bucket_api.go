package bucket

import (
	"encoding/json"
	"gofiel/utils"
	"log"
	"net/http"
)

// TODO: CREATE ENDPOINT FOR GETTING ALL FILENAMES IN THE BUCKET
// MAYBE ADD PAGGINATION (HOW NOT TO QUERY ALL FILES IF WE NEED TO PERFORM os.ReadDir()??)
// ACTUALLY HADRD QUESTION

func RegisterBucketApiEndpoints() {
	http.HandleFunc("/bucket", handleBuckets)
}

func handleBuckets(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getBuckets(w, r)
		return
	case http.MethodPost:
		addBucket(w, r)
		return
	}
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}

func getBuckets(w http.ResponseWriter, r *http.Request) {
	var toReturn []BucketResp = []BucketResp{}
	for _, val := range Buckets {
		toReturn = append(toReturn, val.toBucketResp())
	}
	utils.WriteBasicResp(w, toReturn, 0, "")
}

func addBucket(w http.ResponseWriter, r *http.Request) {
	addBucketReq := AddBucketReq{}

	err := json.NewDecoder(r.Body).Decode(&addBucketReq)

	if err != nil {
		log.Println(err.Error())
		utils.WriteBasicResp(w, nil, 7, "invalid request")
		return
	}

	if addBucketReq.Name == "" {
		log.Println(err.Error())
		utils.WriteBasicResp(w, nil, 7, "invalid invalid bucket-name")
		return
	}

	newBucket, err := createNewBucket(addBucketReq.Name)

	if err != nil {
		// TODO: RollbackAllCreateNewBucketStateModifications(addBucketReq.Name)
		log.Println(err.Error())
		utils.WriteBasicResp(w, nil, 7, "failed while creating the bucket")
		return
	}

	utils.WriteBasicResp(w, newBucket, 0, "")
}
