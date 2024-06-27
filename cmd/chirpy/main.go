package main

import (
	"os"

	"github.com/benskia/Chirpy/internal/chirpyserver"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	jwtSecret := os.Getenv("JWT_SECRET")
	chirpyserver.StartChirpyServer(jwtSecret)
}
