package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
}

func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}
	return db, db.ensureDB()
}

func (db *DB) ensureDB() error {
	_, err := os.Stat(db.path)
	if errors.Is(err, os.ErrNotExist) {
		db.writeDB(DBStructure{
			Chirps: map[int]Chirp{},
			Users:  map[int]User{},
		})
		return nil
	}
	return err
}

func (db *DB) loadDB() (DBStructure, error) {
	dbStructure := DBStructure{}
	f, err := os.ReadFile(db.path)
	if err != nil {
		return dbStructure, err
	}
	err = json.Unmarshal(f, &dbStructure)
	return dbStructure, err
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	dbToWrite, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}
	return os.WriteFile(db.path, dbToWrite, 0644)
}
