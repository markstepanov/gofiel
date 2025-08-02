package utils

import (
	"encoding/json"
	"net/http"
)


type BasicResp struct {
	Code int
	Description string
	Data interface{}
}

func WriteBasicResp(w http.ResponseWriter, body interface{} , errCode int, desc string) {
	w.Header().Add("Content-Type", "application/json")
	resp := BasicResp{
		Code:errCode,
		Description: desc,
		Data: body,
	}

	json, err  := json.Marshal(resp)
	if err != nil {
		return
	}
	w.Write(json)
}


func PostMethod(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		http.HandlerFunc(next).ServeHTTP(w, r)
	})
}


func GetMethod(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet{
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		http.HandlerFunc(next).ServeHTTP(w, r)
	})
}