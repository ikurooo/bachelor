package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", HandleHttp)

	log.Println("Starting standard server on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
