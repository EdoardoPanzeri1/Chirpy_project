package main

import (
	"encoding/json"
	"net/http"
)

func handleLenght(w http.ResponseWriter, r *http.Request) {
	type Chirpy struct {
		Body string `json:"body"`
	}

	var chirp Chirpy
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&chirp)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	if len(chirp.Body) > 140 {
		responseWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	responseWithSuccess(w)
}

func responseWithError(w http.ResponseWriter, code int, message string) {
	response := map[string]string{"error": message}
	jsonResponse, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(jsonResponse)
}

func responseWithSuccess(w http.ResponseWriter) {
	response := map[string]bool{"valid": true}
	jsonResponse, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
