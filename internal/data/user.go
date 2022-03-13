package data

import "time"

// TODO: create user migration
// TODO: Insert
// TODO: GetByEmail
// TODO: GetForToken
// TODO: ValidateUser -> TODO: IsAnon check
// TODO: ValidatePassword -> TODO: SetPass and MatchPass
// TODO: ValidateEmail
// TODO: Custom error for duplicate email

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
