package chirpyserver

import (
	"encoding/json"
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
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	t, _ := template.ParseFiles(templatePath + "metrics.html")
	t.Execute(w, cfg)
}

func (cfg *apiConfig) resetHandle(w http.ResponseWriter, _ *http.Request) {
	cfg.fileserverHits = 0
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Fileserver hit counter has been reset."))
	if err != nil {
		log.Println("Error handling /reset: ", err)
	}
}

func validateChirpHandle(w http.ResponseWriter, r *http.Request) {
	type failResponse struct {
		Error string `json:"error"`
	}
	type successResponse struct {
		Valid bool `json:"valid"`
	}
	type parameters struct {
		Body string `json:"body"`
	}

	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	fail := failResponse{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Println("Error decoding parameters: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		fail.Error = "Something went wrong"
		dat, err := json.Marshal(fail)
		if err != nil {
			log.Println("Error marshalling JSON: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(dat)
		return
	}

	if len(params.Body) > 140 {
		log.Println("Chirp is too long")
		w.WriteHeader(http.StatusBadRequest)
		fail.Error = "Chirp is too long"
		dat, err := json.Marshal(fail)
		if err != nil {
			log.Println("Error marshalling JSON: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(dat)
		return
	}

	success := successResponse{true}
	dat, err := json.Marshal(success)
	if err != nil {
		log.Println("Error marshalling JSON: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(dat)
}

func readinessHandle(w http.ResponseWriter, _ *http.Request) {
	w.Header()["Content-Type"] = []string{"text/plain; charset=utf-8"}
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		log.Println("Error handling /healthz: ", err)
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
	mux.HandleFunc("POST /api/validate_chirp", validateChirpHandle)

	srv := &http.Server{Handler: mux, Addr: localHost}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
