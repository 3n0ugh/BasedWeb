package data

import (
	"github.com/3n0ugh/BasedWeb/internal/validator"
	"time"
)

// TODO: create user migration
// TODO: Insert
// TODO: GetByEmail
// TODO: GetForToken
// TODO: ValidateUser -> TODO: IsAnon check
// TODO: ValidatePassword -> TODO: SetPass and MatchPass
// TODO: ValidateEmail
// TODO: Custom error for duplicate email

var AnonymousUser = &User{}

type User struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	Activated bool      `json:"activated"`
	Version   int       `json:"-"`
}

type password struct {
	plaintext *string
	hash      []byte
}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

func ValidatePassword(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email",
		"must be valid email address")
}

func ValidateUser(v *validator.Validator, u *User) {
	v.Check(u.Name != "", "name", "must be provided")
	v.Check(len(u.Name) <= 500, "name", "must not be more than 500 bytes long")

	ValidateEmail(v, u.Email)

	if u.Password.plaintext != nil {
		ValidatePassword(v, *u.Password.plaintext)
	}

	if u.Password.hash == nil {
		panic("missing password hash for user")
	}
}
