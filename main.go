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

// go tool componentize-go build -o app.wasm
// wasmtime serve -Scli app.wasm
// docker build -t wasm-serverless:latest .
// docker save wasm-serverless:latest -o /tmp/wasm-serverless.tar
// sudo ctr images import /tmp/wasm-serverless.tar
// sudo ctr images ls | grep wasm-serverless -> docker.io/library/wasm-serverless:latest
// sudo ctr run --rm --runtime=io.containerd.wasmtime.v1 docker.io/library/wasm-serverless:latest wasmtest
// sudo ctr containers ls | grep wasmtest
// sudo ctr tasks ls | grep wasmtest
// sudo ctr tasks ls n -> PID = 305706
// sudo nsenter -t 305706 -n ip addr
// sudo nsenter -t 305706 -n ip route
// sudo nsenter -t 305706 -n ip link show lo
// sudo nsenter -t 305706 -n ip link set lo up
// sudo nsenter -t 305706 -n ip addr add 127.0.0.1/8 dev lo
// sudo nsenter -t 305706 -n curl -v http://127.0.0.1:8080/
func main() {}
