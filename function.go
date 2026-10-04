package main

import (
	"net/http"
)

// HandleHttp is your shared business logic
func HandleHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	if r.Method == "GET" && r.URL.Path == "/hello" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello from shared function!\n"))
		return
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Not found\n"))
}
