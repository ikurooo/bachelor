package main

import (
	"fmt"

	"wasm-serverless/function"
)

func main() {
	req := function.Request{
		Method: "GET",
		Path:   "/hello",
	}

	resp := function.Handle(req)

	fmt.Printf("Status: %d\n", resp.Status)
	fmt.Printf("Body: %s", resp.Body)
}
