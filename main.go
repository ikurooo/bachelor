package main

import (
	"log"
	"net/http"
)

// 1. Build the native Docker image using Dockerfile.native
// docker build -f Dockerfile.native -t native-app:latest .

// 2. Save the image to a tarball so K3s can import it
// docker save native-app:latest -o /tmp/native-app.tar

// 3. Import the tarball into K3s containerd
// sudo k3s ctr images import /tmp/native-app.tar
func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("X-Execution-Env", "native")
		HandleHttp(w, r)
	})

	log.Println("Starting standard server on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
