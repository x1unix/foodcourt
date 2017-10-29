package model

import (
	"errors"
	"database/sql"
)

// Superprevileged user
const LEVEL_ADMIN = 0;

// Manager level
const LEVEL_MANAGER = 1;

// Regular user (client)
const LEVEL_USER = 2;

// Application user (`users`)
type User struct {
	ID int
	Email string
	FirstName string
	LastName string
	Password string
	Level int8
	db *sql.DB
}

// Find users by id
func (u *User) FindById(id int) (error, *User) {
	return errors.New("Not Implemented Yet"), nil
}

// Find user
func (u *User) Find(query string) (error, []*User) {
	return errors.New("Not Implemented Yet"), nil
}

// Create user
func (u *User) Create() (error, *User) {
	return errors.New("Not Implemented Yet"), nil
}

// Update user
func (u *User) Update() (error, *User) {
	return errors.New("Not Implemented Yet"), nil
}

// Delete user
func (u *User) Delete() (error, *User) {
	return errors.New("Not Implemented Yet"), nil
}

// Get all users
func (u *User) GetAll() (error, []*User) {
	return errors.New("Not Implemented Yet"), nil
}
