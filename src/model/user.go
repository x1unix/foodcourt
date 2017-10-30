package model

import (
	"errors"
	"database/sql"
)

import (
	sq "github.com/Masterminds/squirrel"
	"fmt"
)

// Superprevileged user
const LEVEL_ADMIN = 0;

// Manager level
const LEVEL_MANAGER = 1;

// Regular user (client)
const LEVEL_USER = 2;

const T_USERS = "users";

// Application user (`users`)
type User struct {
	ID         int
	Email      string
	FirstName  string
	LastName   string
	Password   string
	Level      int8
	Connection *sql.DB
}

// Find users by id
func (u *User) FindById(id int) (error, *User) {
	data, err := sq.Select("*").From(T_USERS).Where(sq.Eq{"id": id}).RunWith(u.Connection).Query();

	if (err == nil) {
		fmt.Println(data)
	}

	return err, nil
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

func (u *User) Dispose() {
	u.Connection.Close();
}
