package auth

// Application user (`users`)
type User struct {
	ID        int      `db:"id" json:"id" msgpack:"id"`
	Email     string   `db:"email" json:"email" msgpack:"email" validate:"required,email"`
	FirstName string   `db:"firstName" json:"firstName" msgpack:"firstName" validate:"required"`
	LastName  string   `db:"lastName" json:"lastName" msgpack:"lastName" validate:"required"`
	Password  string   `db:"password" json:"password,omitempty" msgpack:"-" validate:"required"`
	Level     int     `db:"level" json:"level" msgpack:""`
}
