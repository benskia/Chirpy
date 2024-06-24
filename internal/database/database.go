package database

import (
	"errors"
	"io/fs"
	"log"
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
	err := db.ensureDB()
	if err != nil {
		log.Fatal("Error creating database connection: ", err)
	}
	return db, err
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	return Chirp{
		Body: body,
	}, nil
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
	return DBStructure{}, nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	return nil
}
