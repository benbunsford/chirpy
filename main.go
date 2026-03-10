package main

import (
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(".")))

	port := "8080"
	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	server.ListenAndServe()
}
