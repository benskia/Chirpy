package chirpyserver

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strings"

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

// API/CHIRPS

func getChirps(w http.ResponseWriter, r *http.Request) {
	return
}

func postChirps(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		ID   int    `json:"id"`
		Body string `json:"body"`
	}

	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Println("Error decoding parameters: ", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	if len(params.Body) > 140 {
		msg := "Chirp is too long"
		log.Println(msg)
		respondWithError(w, http.StatusBadRequest, msg)
		return
	}

	cleanedBody := cleanChirp(params.Body)
	payload := map[string]string{"cleaned_body": cleanedBody}
	respondWithJSON(w, http.StatusOK, payload)
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

func cleanChirp(chirp string) string {
	if len(chirp) == 0 {
		return chirp
	}
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleanedWords := []string{}
	splitChirp := strings.Split(chirp, " ")
	for _, word := range splitChirp {
		normalizedWord := strings.ToLower(word)
		if _, ok := badWords[normalizedWord]; ok {
			cleanedWords = append(cleanedWords, "****")
			continue
		}
		cleanedWords = append(cleanedWords, word)
	}
	return strings.Join(cleanedWords, " ")
}

// API/HEALTHZ

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
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandle)
	mux.HandleFunc("GET /api/reset", apiCfg.resetHandle)
	mux.HandleFunc("GET /api/healthz", readinessHandle)
	mux.HandleFunc("GET /api/chirps", getChirps)
	mux.HandleFunc("POST /api/chirps", postChirps)

	srv := &http.Server{Handler: mux, Addr: localHost}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
