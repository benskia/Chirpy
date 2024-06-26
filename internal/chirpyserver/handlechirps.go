package chirpyserver

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func (cfg *apiConfig) getChirp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	targetID, err := strconv.Atoi(r.PathValue("chirpID"))
	if err != nil {
		log.Println("Error converting ID to integer for getChirp: ", err)
	}
	chirps, err := cfg.db.GetChirps()
	if err != nil {
		log.Println("Error retrieving chirps for getChirp: ", err)
	}
	for _, chirp := range chirps {
		if chirp.ID == targetID {
			respondWithJSON(w, http.StatusOK, chirp)
			return
		}
	}
	respondWithError(w, http.StatusNotFound, "Chirp not found")
}

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	chirps, err := cfg.db.GetChirps()
	if err != nil {
		log.Println("Error retrieving chirps for getChirps: ", err)
		return
	}
	encodedChirps, err := json.Marshal(chirps)
	if err != nil {
		log.Println("Error marshalling chirps for getChirps: ", err)
		return
	}
	_, err = w.Write(encodedChirps)
	if err != nil {
		log.Println("Error writing chirps for getChirps: ", err)
	}
}

func (cfg *apiConfig) postChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
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
	chirp, err := cfg.db.CreateChirp(cleanedBody)
	if err != nil {
		log.Println("Error creating chirp: ", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, chirp)
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
