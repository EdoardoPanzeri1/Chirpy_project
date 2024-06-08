package main

import (
	"fmt"
	"net/http"
)

func main() {
	// Create a new ServeMux
	mux := http.NewServeMux()

	// Create a new http.Server and assign the mux to it
	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}
	// Handle the root path
	mux.Handle("/", http.FileServer(http.Dir(".")))

	// Start the server and check for errors
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}

}
