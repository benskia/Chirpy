package chirpyserver

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/benskia/Chirpy/internal/database"
	"golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) postUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Println("Error decoding parameters: ", err)
		respondWithError(w, http.StatusBadRequest, internalErrorMsg)
		return
	}
	user, err := cfg.db.CreateUser(params.Email, params.Password)
	if err != nil {
		log.Println("Error creating user: ", err)
		respondWithError(w, http.StatusInternalServerError, internalErrorMsg)
		return
	}
	respondWithJSON(w, http.StatusCreated, user)
}

func (cfg *apiConfig) loginUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Println("Error decoding parameters: ", err)
		respondWithError(w, http.StatusBadRequest, internalErrorMsg)
		return
	}
	users, err := cfg.db.GetUsers()
	if err != nil {
		log.Println("Error retrieving users: ", err)
		respondWithError(w, http.StatusInternalServerError, internalErrorMsg)
		return
	}
	for _, user := range users {
		err := bcrypt.CompareHashAndPassword(user.Password, []byte(params.Password))
		if err == nil {
			userResponse := database.UserResponse{
				ID:    user.ID,
				Email: user.Email,
			}
			respondWithJSON(w, http.StatusOK, userResponse)
			return
		}
	}
	respondWithError(w, http.StatusUnauthorized, "Invalid username/password")
}
