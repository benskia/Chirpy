#!/bin/bash

rm database.json

go build -o ./bin/chirpy ./cmd/chirpy/main.go && ./bin/chirpy
