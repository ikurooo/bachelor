package main

import (
	"net/http"

	"go.bytecodealliance.org/pkg/wasihttp"
)

func init() {
	wasihttp.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello from my WASI HTTP function!\n"))
	})
}

func main() {}
