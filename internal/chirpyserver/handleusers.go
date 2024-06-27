package chirpyserver

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/benskia/Chirpy/internal/database"
	"github.com/golang-jwt/jwt/v5"
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

func (cfg *apiConfig) putUser(w http.ResponseWriter, r *http.Request) {
	return
}

func (cfg *apiConfig) loginUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email          string        `json:"email"`
		Password       string        `json:"password"`
		ExpirationSecs time.Duration `json:"expires_in_seconds"`
	}
	defaultExpiration := time.Second * 86400 // seconds in a day
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	params := parameters{"", "", defaultExpiration}
	err := decoder.Decode(&params)
	if err != nil {
		log.Println("Error decoding parameters: ", err)
		respondWithError(w, http.StatusBadRequest, internalErrorMsg)
		return
	}
	if params.ExpirationSecs > defaultExpiration {
		params.ExpirationSecs = defaultExpiration
	}
	users, err := cfg.db.GetUsers()
	if err != nil {
		log.Println("Error retrieving users: ", err)
		respondWithError(w, http.StatusInternalServerError, internalErrorMsg)
		return
	}
	start := time.Now()
	end := start.Add(params.ExpirationSecs)
	for _, user := range users {
		err := bcrypt.CompareHashAndPassword(user.Password, []byte(params.Password))
		if err == nil {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
				Issuer:    "chirpy",
				IssuedAt:  jwt.NewNumericDate(start),
				ExpiresAt: jwt.NewNumericDate(end),
				Subject:   strconv.Itoa(user.ID),
			})
			signedToken, err := token.SignedString([]byte(cfg.secret))
			if err != nil {
				log.Println("Error signing token: ", err)
				respondWithError(w, http.StatusInternalServerError, internalErrorMsg)
				return
			}
			userResponse := database.UserResponse{
				ID:    user.ID,
				Email: user.Email,
				Token: signedToken,
			}
			respondWithJSON(w, http.StatusOK, userResponse)
			return
		}
	}
	respondWithError(w, http.StatusUnauthorized, "Invalid username/password")
}
