package handler

import (
	"fmt"
	"net/http"
)

func HandleHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	environment := r.Header.Get("X-Execution-Env")

	if r.Method == "GET" && r.URL.Path == "/hello" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Hello from %s function!\n", environment)))
		return
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Not found\n"))
}
