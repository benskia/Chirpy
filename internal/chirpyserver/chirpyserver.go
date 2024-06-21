package chirpyserver

import (
	"html/template"
	"log"
	"net/http"
)

const (
	localHost    string = "localhost:8080"
	websitePath  string = "./website/"
	templatePath string = "./web/"
)

type apiConfig struct {
	fileserverHits int
}

func NewApiConfig() *apiConfig {
	return &apiConfig{0}
}

func (cfg *apiConfig) GetHits() int {
	return cfg.fileserverHits
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) metricsHandle(w http.ResponseWriter, _ *http.Request) {
	w.Header()["Content-Type"] = []string{"text/html; charset=utf-8"}
	w.WriteHeader(http.StatusOK)
	t, _ := template.ParseFiles(templatePath + "metrics.html")
	t.Execute(w, cfg)
}

func (cfg *apiConfig) resetHandle(w http.ResponseWriter, _ *http.Request) {
	cfg.fileserverHits = 0
	w.Header()["Content-Type"] = []string{"text/plain; charset=utf-8"}
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Fileserver hit counter has been reset."))
	if err != nil {
		log.Println("/reset failed to write body data")
	}
}

func StartChirpyServer() {
	apiCfg := NewApiConfig()
	handler := http.StripPrefix("/app", http.FileServer(http.Dir(websitePath)))
	mux := http.NewServeMux()
	mux.Handle("GET /app/*", apiCfg.middlewareMetricsInc(handler))
	mux.HandleFunc("GET /api/healthz", readinessHandle)
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandle)
	mux.HandleFunc("GET /api/reset", apiCfg.resetHandle)

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
