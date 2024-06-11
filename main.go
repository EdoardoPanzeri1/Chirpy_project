package main

import (
	"fmt"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
}

func main() {
	apiCfg := &apiConfig{fileserverHits: 0}

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Wrap the handler with middleware
	fileServerHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	mux.Handle("/app/", fileServerHandler)

	// Handle readiness endpoint
	mux.HandleFunc("GET /api/healthz", handleReadiness)

	// Handle the count
	mux.HandleFunc("GET /admin/metrics", apiCfg.handleMetrics)

	// Handle the reset
	mux.HandleFunc("GET /api/reset", apiCfg.handleReset)

	// Handle the lenght of the messages
	mux.HandleFunc("POST /api/validate_chirp", handleLenght)

	// Create a new http.Server and assign the mux to it
	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}

	// Start the server and check for errors
	fmt.Println("Starting server on :8080")
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Error starting server: ", err)
	}
}
