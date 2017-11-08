package model

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"crypto/md5"
	"encoding/hex"
)

// Superprevileged user
const LEVEL_ADMIN = 0;

// Manager level
const LEVEL_MANAGER = 1;

// Regular user (client)
const LEVEL_USER = 2;

const T_USERS = "users";

// User id column key
const UserId = "id";

// User
const UserEmail = "email";

// Application user (`users`)
type User struct {
	ID        int      `db:"id" json:"id" msgpack:"id"`
	Email     string   `db:"email" json:"email" msgpack:"email" validate:"required,email"`
	FirstName string   `db:"firstName" json:"firstName" msgpack:"firstName" validate:"required"`
	LastName  string   `db:"lastName" json:"lastName" msgpack:"lastName" validate:"required"`
	Password  string   `db:"password" json:"password,omitempty" msgpack:"-" validate:"required"`
	Level     int     `db:"level" json:"level" msgpack:""`
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
func (u *User) Find(query sq.Eq) (error, *User) {
	q, args, _ := sq.Select("*").From(T_USERS).Where(query).ToSql()
	user := &User{}

	err := u.DB.Get(user, q, args...)

	return err, user
}

func (u *User) FindAll(query sq.Eq) (error, *[]User) {
	q, args, _ := sq.Select("*").From(T_USERS).Where(query).ToSql()
	users := &[]User{}

	err := u.DB.Get(users, q, args...)

	return err, users
}

// Check if user id exists
func (u *User) IdExists() (error, bool) {
	q, _, _ := sq.Select("COUNT(*)").From(T_USERS).Where(sq.Eq{"id": u.ID}).ToSql()

	count := 0
	err := u.DB.Get(&count, q, u.ID)
	ifExists := count > 0

	return err, ifExists
}

// check if user exists
func (u *User) Exists() (error, bool) {
	q, _, _ := sq.Select("COUNT(*)").From(T_USERS).Where(sq.Eq{"email": u.Email}).ToSql()

	count := 0
	err := u.DB.Get(&count, q, u.Email)
	ifExists := count > 0

	return err, ifExists
}

// Create user
func (u *User) Create() error {
	q, _, _ := sq.Insert(T_USERS).
	Columns("Email", "FirstName", "LastName", "Password", "Level").
	Values(u.Email, u.FirstName, u.LastName, u.Password, u.Level).
	ToSql()

	_, err := u.DB.Exec(q, u.Email, u.FirstName, u.LastName, u.Password, u.Level)

	return err
}

// Update user
func (u *User) Update() error {
	builder := sq.Update(T_USERS).Where(sq.Eq{"id": u.ID})

	// TODO: add method to autofill UPDATE query with non-empty fields
	if (u.FirstName != "") {
		builder = builder.Set("firstName", u.FirstName)
	}

	if (u.LastName != "") {
		builder = builder.Set("lastName", u.LastName)
	}

	if (u.Password != "") {
		builder = builder.Set("password", u.Password)
	}

	_, err := builder.RunWith(u.DB.DB).Query()

	return err
}

// Delete user
func (u *User) Delete() error {
	q, _, _ := sq.Delete(T_USERS).
		Where(sq.Eq{"id": u.ID}).
		ToSql()

	_, err := u.DB.Exec(q, u.ID)

	return err
}

func (u *User) HashPassword(pass string) string {
	hasher := md5.New()
	hasher.Write([]byte(pass))
	return hex.EncodeToString(hasher.Sum(nil))
}

// Validate user credentials
func (u *User) MatchCredentials(pass string) bool {
	passwd := u.HashPassword(pass)

	return passwd == u.Password
}

// Get ID of current user by defined email
func (u *User) GetId() error {
	q, _, _ := sq.Select("id").
		From(T_USERS).
		Where(sq.Eq{"email": u.Email}).
		ToSql()

	var uid int
	err := u.DB.Get(&uid, q, u.Email)

	if (err == nil) {
		u.ID = uid
	}

	return err
}

func (u *User) Dispose() {
	defer u.DB.Close();
}

func Users(con *sqlx.DB) *User {
	u := User{DB: con}
	return &u
}