package bucket

import (
	"gofiel/utils"
	"net/http"
)

// TODO: CREATE ENDPOINT FOR GETTING ALL FILENAMES IN THE BUCKET
// MAYBE ADD PAGGINATION (HOW NOT TO QUERY ALL FILES IF WE NEED TO PERFORM os.ReadDir()??)
// ACTUALLY HADRD QUESTION

func RegisterStorageApiEndpoints() {
	http.HandleFunc("/bucket", handleBuckets)
}

func handleBuckets(w http.ResponseWriter, r *http.Request) {
	var toReturn []BucketResp = []BucketResp{}
	for _, val := range Bukets {
		toReturn = append(toReturn, val.toBucketResp())
	}
	utils.WriteBasicResp(w, toReturn, 0, "")
}
