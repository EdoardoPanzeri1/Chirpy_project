package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handleLenght(w http.ResponseWriter, r *http.Request) {
	type Chirpy struct {
		Body string `json:"body"`
	}

	banWords := []string{"kerfuffle", "sharbert", "fornax"}

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

	words := strings.Split(chirp.Body, " ")

	for i, word := range words {
		for _, banWord := range banWords {
			if strings.ToLower(word) == banWord {
				words[i] = "****"
				break
			}
		}
	}

	cleanedBody := strings.Join(words, " ")
	responseWithJSON(w, 200, map[string]string{"cleaned_body": cleanedBody})
}

func responseWithError(w http.ResponseWriter, code int, message string) {
	response := map[string]string{"error": message}
	jsonResponse, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(jsonResponse)
}

func responseWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	jsonResponse, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
