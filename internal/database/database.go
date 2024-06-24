package database

import (
	"log"
	"sort"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

func NewDB(path string) (*DB, error) {
	return &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}, nil
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	return Chirp{
		Body: body,
	}, nil
}

func (db *DB) GetChirps(body string) ([]Chirp, error) {
	loadedDB, err := db.loadDB()
	if err != nil {
		log.Println("Error loading database: ", err)
	}
	chirps := []Chirp{}
	for index, chirp := range loadedDB.Chirps {
		chirps = append(chirps, Chirp{index, chirp.Body})
	}
	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})
	return chirps, nil
}

func (db *DB) ensureDB() error {
	return nil
}

func (db *DB) loadDB() (DBStructure, error) {
	return DBStructure{}, nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	return nil
}
