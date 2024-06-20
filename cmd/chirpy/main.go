package main

import (
	"log"
	"net/http"
)

func main() {
	const localHost string = "localhost:8080"
	mux := http.NewServeMux()
	srv := http.Server{Handler: mux, Addr: localHost}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
