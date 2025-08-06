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
	for _, val := range Bukets {
		toReturn = append(toReturn, val.toBucketResp())
	}
	utils.WriteBasicResp(w, toReturn, 0, "")
}

func addBucket(w http.ResponseWriter, r *http.Request) {
	addBucketResp := AddBucketResp{}

	err := json.NewDecoder(r.Body).Decode(&addBucketResp)

	if err != nil {
		log.Println(err.Error())
		utils.WriteBasicResp(w, nil, 7, "invalid request")
		return
	}

	if addBucketResp.Name == "" {
		utils.WriteBasicResp(w, nil, 7, "invalid invalid bucket-name")
		return
	}

	utils.WriteBasicResp(w, Bucket{Id: 6969, Name: addBucketResp.Name}, 0, "")

}
