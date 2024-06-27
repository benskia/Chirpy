package chirpyserver

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/benskia/Chirpy/internal/database"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) postUser(w http.ResponseWriter, r *http.Request) {
	type response struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
	}
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
	respondWithJSON(w, http.StatusCreated, response{
		ID:    user.ID,
		Email: user.Email,
	})
}

func (cfg *apiConfig) putUser(w http.ResponseWriter, r *http.Request) {
	type response struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
	}
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	authHeader := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	claims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(authHeader, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.secret), nil
	})
	if err != nil {
		log.Println("Error parsing claims: ", err)
		respondWithError(w, http.StatusUnauthorized, "Invalid authorization token")
		return
	}
	subject, err := token.Claims.GetSubject()
	if err != nil {
		log.Println("Error getting subect from claims: ", err)
		respondWithError(w, http.StatusInternalServerError, internalErrorMsg)
		return
	}
	id, err := strconv.Atoi(subject)
	if err != nil {
		log.Println("Error converting claim subject to int ID: ", err)
		respondWithError(w, http.StatusInternalServerError, internalErrorMsg)
		return
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		log.Println("Error decoding parameters: ", err)
		respondWithError(w, http.StatusBadRequest, internalErrorMsg)
		return
	}
	newPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	cfg.db.UpdateUser(database.User{
		ID:       id,
		Email:    params.Email,
		Password: newPassword,
	})
	respondWithJSON(w, http.StatusOK, response{
		ID:    id,
		Email: params.Email,
	})
}

func (cfg *apiConfig) loginUser(w http.ResponseWriter, r *http.Request) {
	type response struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
		Token string `json:"token"`
	}
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
			respondWithJSON(w, http.StatusOK, response{
				ID:    user.ID,
				Email: user.Email,
				Token: signedToken,
			})
			return
		}
	}
	respondWithError(w, http.StatusUnauthorized, "Invalid username/password")
}
