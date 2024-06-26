#!/bin/bash

curl -v http://localhost:8080/app/

curl -v --header "Content-Type: application/json" \
  --request POST \
  --data '{"body":"I hear Mastodon is better than Chirpy. sharbert I need to migrate"}' \
  http://localhost:8080/api/chirps

curl -v --header "Content-Type: application/json" \
  --request GET \
  http://localhost:8080/api/chirps

curl -v --header "Content-Type: application/json" \
  --request POST \
  --data '{"email":"userOne@example.com"}' \
  http://localhost:8080/api/users

curl -v --header "Content-Type: application/json" \
  --request POST \
  --data '{"email":"userTwo@example.com"}' \
  http://localhost:8080/api/users

