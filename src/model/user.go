package model

import (
	"errors"
)

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
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
	ID        int      `db:"id" json:"id"`
	Email     string   `db:"email" json:"email" validate:"required,email"`
	FirstName string   `db:"firstName" json:"firstName" validate:"required"`
	LastName  string   `db:"lastName" json:"lastName" validate:"required"`
	Password  string   `db:"password" json:"-" validate:"required"`
	Level     int8     `db:"level" json:"level"`
	DB        *sqlx.DB `json:"-"`
}

// Find users by id
func (u *User) FindById(id string) (error, *User) {
	q, _, _ := sq.Select("*").From(T_USERS).Where(sq.Eq{"id": id}).ToSql()
	user := User{}

	err := u.DB.Get(&user, q, id)

	return err, &user
}

// Get all users
func (u *User) GetAll() (error, *[]User) {
	q, _, _ := sq.Select("*").From(T_USERS).ToSql()

	users := []User{}
	err := u.DB.Select(&users, q)

	return err, &users
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

func (u *User) Dispose() {
	defer u.DB.Close();
}

func Users(con *sqlx.DB) *User {
	u := User{DB: con}
	return &u
}