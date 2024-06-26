package chirpyserver

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *apiConfig) postUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
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
