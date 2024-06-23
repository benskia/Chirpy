#!/bin/bash

curl -v http://localhost:8080/app/

curl -v --header "Content-Type: application/json" \
  --request POST \
  --data '{"body":"A short chirp."}' \
  http://localhost:8080/api/validate_chirp

curl -v --header "Content-Type: application/json" \
  --request POST \
  --data '{"body":"This is a longer chirp. Chirps cannot contain more than 140 characters, so we'\''re gonna keep typing until we have at least 141 characters. This should do it."}' \
  http://localhost:8080/api/validate_chirp

curl -v --header "Content-Type: application/json" \
  --request POST \
  --data '{"body":"I hear Mastodon is better than Chirpy. sharbert I need to migrate"}' \
  http://localhost:8080/api/validate_chirp
