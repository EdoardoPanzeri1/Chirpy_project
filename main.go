package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/EdoardoPanzeri1/Chirpy_project/database"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	apiCfg := apiConfig{
		fileserverHits: 0,
		DB:             db,
	}

	mux := http.NewServeMux()
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/", fsHandler)

	mux.HandleFunc("/api/healthz", handlerReadiness)
	mux.HandleFunc("/api/reset", apiCfg.handlerReset)
	mux.HandleFunc("/api/chirps", apiCfg.handlerChirps)
	mux.HandleFunc("/api/chirps/", apiCfg.handlerChirpsGet)
	mux.HandleFunc("/api/users", apiCfg.handlerUserCreate)

	mux.HandleFunc("/admin/metrics", apiCfg.handlerMetrics)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

func (cfg *apiConfig) handlerChirps(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		cfg.handlerChirpsCreate(w, r)
	case http.MethodGet:
		cfg.handlerChirpsRetrieve(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

func (cfg *apiConfig) handlerUserCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	user, err := cfg.DB.CreateUser(params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	respondWithJSON(w, http.StatusCreated, User{
		ID:    user.ID,
		Email: user.Email,
	})
}
