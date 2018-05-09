package auth

const NilUserLevel = -1

// Application user (`users`)
type User struct {
	ID        int    `db:"id" json:"id" msgpack:"id"`
	Email     string `db:"email" json:"email" msgpack:"email" validate:"required,email"`
	FirstName string `db:"firstName" json:"firstName" msgpack:"firstName" validate:"required"`
	LastName  string `db:"lastName" json:"lastName" msgpack:"lastName" validate:"required"`
	Password  string `db:"password" json:"password,omitempty" msgpack:"-" validate:"required"`
	Level     int    `db:"level" json:"level" msgpack:""`
}

// Checks if the user belongs to specific group
func (u *User) IsInGroup(group int) bool {
	return u.Level == group
}

func (u *User) IsAdmin() bool {
	return u.IsInGroup(LEVEL_ADMIN)
}

func (u *User) IsManager() bool {
	return u.IsInGroup(LEVEL_MANAGER)
}

func (u *User) IsClient() bool {
	return u.IsInGroup(LEVEL_USER)
}

// GetUserBoilerplate returns user structure boilerplate with uninitialized values.
// Necessary for user edit action (check if user level was changed, etc), because
// only changed values are passed.
func NewUser() (u User) {
	u.Level = NilUserLevel
	return u
}