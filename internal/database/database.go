package database

import "sync"

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

type Chirp struct {
	Body string `json:"body"`
}

func NewDB(path string) (*DB, error) {
	return nil, nil
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	return Chirp{}, nil
}

func (db *DB) GetChirps(body string) ([]Chirp, error) {
	return []Chirp{}, nil
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
