//go:build wasip2

package main

import (
	"fmt"
	"net/http"

	"go.bytecodealliance.org/pkg/wasihttp"
)

func main() {
	wasihttp.Handle(http.HandlerFunc(HandleHttp))
}

func HandleHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("X-Execution-Env", "wasm")

	if r.Method == "GET" && r.URL.Path == "/hello" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Hello from %s function!\n", r.Header.Get("X-Execution-Env"))))
		return
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Not found\n"))
}

// Compile the WebAssembly component using componentize-go
// go tool componentize-go build -o app.wasm

// Build the Wasm Docker image (using your FROM scratch Dockerfile)
// docker build -f Dockerfile -t wasm-serverless:latest .

// Save the Wasm image to a tarball for K3s import
// docker save wasm-serverless:latest -o /tmp/wasm-serverless.tar

// sudo k3s ctr images import /tmp/wasm-serverless.tar
