package main

import (
	"net/http"

	"wasm-serverless/pkg/handler"

	"go.bytecodealliance.org/pkg/wasihttp"
)

func init() {
	wasihttp.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("X-Execution-Env", "wasm")
		handler.HandleHttp(w, r)
	})
}

// Compile the WebAssembly component using componentize-go
// go tool componentize-go build -o app.wasm

// Build the Wasm Docker image (using your FROM scratch Dockerfile)
// docker build -f Dockerfile -t wasm-serverless:latest .

// Save the Wasm image to a tarball for K3s import
// docker save wasm-serverless:latest -o /tmp/wasm-serverless.tar

// sudo k3s ctr images import /tmp/wasm-serverless.tar
func main() {}
