package utils

import (
	"encoding/json"
	"net/http"
)

type BasicResp struct {
	Code        int
	Description string
	Data        any
}

func WriteBasicResp(w http.ResponseWriter, body any, errCode int, desc string) {
	w.Header().Add("Content-Type", "application/json")
	resp := BasicResp{
		Code:        errCode,
		Description: desc,
		Data:        body,
	}

	json, err := json.Marshal(resp)
	if err != nil {
		return
	}
	w.Write(json)
}
