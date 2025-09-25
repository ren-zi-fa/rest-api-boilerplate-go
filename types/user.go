package types

import "time"

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

type RegisterUserPayload struct {
	Username  string    `json:"username" validate:"required,min=3"`
	Password  string    `json:"password" validate:"required,min=6"`
	Email     string    `json:"email" validate:"required,email"`
	CreatedAt time.Time `json:"createdAt"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUsers() ([]*User, error)
	CreateUser(user *User) (int64, error)
}
