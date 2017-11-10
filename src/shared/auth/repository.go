package auth

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

// User id column key
const UserId = "id";

// User
const UserEmail = "email";


// Find users by id
func FindById(db *sqlx.DB, id string) (error, *User) {
	q, _, _ := sq.Select("*").From(T_USERS).Where(sq.Eq{"id": id}).ToSql()
	user := User{}

	err := db.Get(&user, q, id)

	return err, &user
}

// Get all users
func GetAll(db *sqlx.DB, ) (error, *[]User) {
	q, _, _ := sq.Select("*").From(T_USERS).ToSql()

	users := []User{}
	err := db.Select(&users, q)

	return err, &users
}

func UserExists(db *sqlx.DB, query sq.Eq) (error, bool) {
	q, args, _ := sq.Select("COUNT(*)").From(T_USERS).Where(query).ToSql()
	var count int

	err := db.Get(&count, q, args...)

	return err, count > 0
}

// Find user
func Find(db *sqlx.DB, query sq.Eq) (error, *User) {
	q, args, _ := sq.Select("*").From(T_USERS).Where(query).ToSql()
	user := &User{}

	err := db.Get(user, q, args...)

	return err, user
}

func (u *User) FindAll(db *sqlx.DB, query sq.Eq) (error, *[]User) {
	q, args, _ := sq.Select("*").From(T_USERS).Where(query).ToSql()
	users := &[]User{}

	err := db.Get(users, q, args...)

	return err, users
}

// Check if user id exists
func IdExists(db *sqlx.DB, id int) (error, bool) {
	return UserExists(db, sq.Eq{"id": id})
}

// check if user exists
func MailExists(db *sqlx.DB, email string) (error, bool) {
	return UserExists(db, sq.Eq{"email": email})
}

// Create user
func AddUser(db *sqlx.DB, u *User) error {
	password := Password(u.Password)

	_, err := sq.Insert(T_USERS).
		Columns("email", "firstName", "lastName", "password", "level").
		Values(u.Email, u.FirstName, u.LastName,  password, u.Level).
		RunWith(db.DB).
		Exec()

	return err
}

// Update user
func UpdateUser(db *sqlx.DB, u *User) error {
	builder := sq.Update(T_USERS).Where(sq.Eq{"id": u.ID})

	// TODO: add method to autofill UPDATE query with non-empty fields
	if u.FirstName != "" {
		builder = builder.Set("firstName", u.FirstName)
	}

	if u.LastName != "" {
		builder = builder.Set("lastName", u.LastName)
	}

	if u.Password != "" {
		password := Password(u.Password)
		builder = builder.Set("password", password)
	}

	_, err := builder.RunWith(db.DB).Query()

	return err
}

// Delete user
func Delete(db *sqlx.DB, id int) error {
	q, args, _ := sq.Delete(T_USERS).
		Where(sq.Eq{"id": id}).
		ToSql()

	_, err := db.Exec(q, args...)

	return err
}

// Get ID of current user by defined email
func GetIdForEmail(db *sqlx.DB, mail string) (error, int) {
	q, args, _ := sq.Select("id").
		From(T_USERS).
		Where(sq.Eq{"email": mail}).
		ToSql()

	var uid int
	err := db.Get(&uid, q, args...)

	return err, uid
}
