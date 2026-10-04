//go:build wasip2

package main

import (
	"go.bytecodealliance.org/pkg/wasihttp"
)

func init() {
	wasihttp.HandleFunc(HandleHttp)
}

func main() {}
