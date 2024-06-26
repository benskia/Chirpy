package database

import (
	"sort"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password []byte `json:"password"`
}

type UserResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

func (db *DB) CreateUser(email, password string) (UserResponse, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return UserResponse{}, err
	}
	index := len(dbStructure.Users) + 1
	bcryptedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return UserResponse{}, err
	}
	newUser := User{
		ID:       index,
		Email:    email,
		Password: bcryptedPass,
	}
	dbStructure.Users[index] = newUser
	err = db.writeDB(dbStructure)
	if err != nil {
		return UserResponse{}, err
	}
	return UserResponse{
		ID:    index,
		Email: email,
	}, nil
}

func (db *DB) GetUsers() ([]User, error) {
	loadedDB, err := db.loadDB()
	if err != nil {
		return []User{}, err
	}
	users := []User{}
	for index, user := range loadedDB.Users {
		users = append(users, User{index, user.Email, user.Password})
	}
	sort.Slice(users, func(i, j int) bool {
		return users[i].ID < users[j].ID
	})
	return users, nil
}
