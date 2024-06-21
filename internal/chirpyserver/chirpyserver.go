package chirpyserver

import (
	"log"
	"net/http"
	"strconv"
)

type apiConfig struct {
	fileserverHits int
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	cfg.fileserverHits++
	return next
}

func (cfg *apiConfig) metricsHandle(w http.ResponseWriter, _ *http.Request) {
	w.Header()["Content-Type"] = []string{"text/plain; charset=utf-8"}
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(strconv.Itoa(cfg.fileserverHits)))
	if err != nil {
		log.Println("/metrics failed to write body data")
	}
}

func (cfg *apiConfig) resetHandle(w http.ResponseWriter, _ *http.Request) {
	cfg.fileserverHits = 0
	w.Header()["Content-Type"] = []string{"text/plain; charset=utf-8"}
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(strconv.Itoa(cfg.fileserverHits)))
	if err != nil {
		log.Println("/reset failed to write body data")
	}
}

func StartChirpyServer() {
	const (
		localHost   string = "localhost:8080"
		websitePath string = "./website/"
	)

	mux := http.NewServeMux()
	apiCfg := apiConfig{}
	handler := http.FileServer(http.Dir(websitePath))
	mux.Handle("/app/*", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(handler)))
	mux.HandleFunc("/healthz", readinessHandle)
	mux.HandleFunc("/metrics", apiCfg.metricsHandle)
	mux.HandleFunc("/reset", apiCfg.resetHandle)

	srv := &http.Server{Handler: mux, Addr: localHost}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func readinessHandle(w http.ResponseWriter, _ *http.Request) {
	w.Header()["Content-Type"] = []string{"text/plain; charset=utf-8"}
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		log.Println("/healthz failed to write body data")
	}
}
