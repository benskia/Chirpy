package main

import "net/http"

func main() {
	mux := http.NewServeMux()
	server := http.Server{Handler: mux}
	server.ListenAndServe()
	return
}
