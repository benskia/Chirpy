package database

import (
	"errors"
	"sort"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password []byte `json:"password"`
	Token    string `json:"token"`
}

func (db *DB) CreateUser(email, password string) (User, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	index := len(dbStructure.Users) + 1
	bcryptedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}
	newUser := User{
		ID:       index,
		Email:    email,
		Password: bcryptedPass,
	}
	dbStructure.Users[index] = newUser
	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}
	return User{
		ID:    index,
		Email: email,
	}, nil
}

func (db *DB) GetUsers() ([]User, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()
	loadedDB, err := db.loadDB()
	if err != nil {
		return []User{}, err
	}
	users := []User{}
	for index, user := range loadedDB.Users {
		users = append(users, User{index, user.Email, user.Password, ""})
	}
	sort.Slice(users, func(i, j int) bool {
		return users[i].ID < users[j].ID
	})
	return users, nil
}

func (db *DB) UpdateUser(user User) error {
	db.mux.Lock()
	defer db.mux.Unlock()
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}
	if _, ok := dbStructure.Users[user.ID]; !ok {
		return errors.New("Failed to update: User not found.")
	}
	dbStructure.Users[user.ID] = user
	db.writeDB(dbStructure)
	return nil
}
