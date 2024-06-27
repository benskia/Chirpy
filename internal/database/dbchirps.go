package database

import "sort"

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	index := len(dbStructure.Chirps) + 1
	newChirp := Chirp{
		ID:   index,
		Body: body,
	}
	dbStructure.Chirps[index] = newChirp
	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}
	return newChirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()
	loadedDB, err := db.loadDB()
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
