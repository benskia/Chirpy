package chirpyserver

import (
	"log"
	"net/http"
)

const (
	localHost   string = "localhost:8080"
	contentPath string = "./website/content/"
	assetsPath  string = "./website/assets/"
)

func StartChirpyServer() {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir(contentPath)))
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(assetsPath))))

	srv := http.Server{Handler: mux, Addr: localHost}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
