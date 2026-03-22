package main

import (
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir("."))))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	port := "8080"
	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	server.ListenAndServe()
}
