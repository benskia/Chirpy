package chirpyserver

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/benskia/Chirpy/internal/database"
)

const (
	localHost    string = "localhost:8080"
	websitePath  string = "./website/"
	templatePath string = "./web/"
)

type apiConfig struct {
	fileserverHits int
	chirpIdCounter int
	db             *database.DB
}

func NewApiConfig() *apiConfig {
	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}
	return &apiConfig{0, 1, db}
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

func readinessHandle(w http.ResponseWriter, _ *http.Request) {
	w.Header()["Content-Type"] = []string{"text/plain; charset=utf-8"}
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		log.Println("Error handling /healthz: ", err)
	}
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type failResponse struct {
		Error string `json:"error"`
	}
	response := failResponse{msg}
	dat, err := json.Marshal(response)
	if err != nil {
		log.Println("Error marshalling JSON: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Println("Error encoding payload: ", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
	}
}

func StartChirpyServer() {
	apiCfg := NewApiConfig()
	handler := http.StripPrefix("/app", http.FileServer(http.Dir(websitePath)))
	mux := http.NewServeMux()
	mux.Handle("GET /app/*", apiCfg.middlewareMetricsInc(handler))
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandle)
	mux.HandleFunc("GET /api/reset", apiCfg.resetHandle)
	mux.HandleFunc("GET /api/healthz", readinessHandle)
	mux.HandleFunc("GET /api/chirps", apiCfg.getChirps)
	mux.HandleFunc("POST /api/chirps", apiCfg.postChirp)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.getChirp)
	mux.HandleFunc("POST /api/users", apiCfg.postUser)
	mux.HandleFunc("POST /api/login", apiCfg.loginUser)

	srv := &http.Server{Handler: mux, Addr: localHost}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
