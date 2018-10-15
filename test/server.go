package main

import (
	"io"
	"net/http"
)

func main() {
	// Hello world, the web server

	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Hello, world!\n")
	}

	http.HandleFunc("/hello", helloHandler)
	go http.ListenAndServe(":", nil)
	http.ListenAndServe("127.0.0.1:", nil)
}
