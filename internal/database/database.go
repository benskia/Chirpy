package database

import (
	"encoding/json"
	"os"
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
	db := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}
	return db, db.ensureDB()
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	index := db.nextIndex(dbStructure)
	newChirp := Chirp{
		ID:   index,
		Body: body,
	}
	dbStructure.Chirps[index] = newChirp
	err = db.writeDB(dbStructure)
	return newChirp, err
}

func (db *DB) GetChirps(body string) ([]Chirp, error) {
	db.mux.RLock()
	loadedDB, err := db.loadDB()
	db.mux.RUnlock()
	if err != nil {
		return []Chirp{}, err
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
	f, err := os.OpenFile(db.path, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	return f.Close()
}

func (db *DB) loadDB() (DBStructure, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()
	dbStructure := DBStructure{}
	f, err := os.ReadFile(db.path)
	if err != nil {
		return dbStructure, err
	}
	err = json.Unmarshal(f, &dbStructure)
	return dbStructure, err
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()
	dbToWrite, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}
	return os.WriteFile(db.path, dbToWrite, 0644)
}

func (db *DB) nextIndex(dbStructure DBStructure) int {
	next := 0
	for i := range dbStructure.Chirps {
		if i > next {
			next = i
		}
	}
	return next + 1
}
