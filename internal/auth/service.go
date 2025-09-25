package auth

import (
	"time"

	"github.com/ren-zi-fa/rest-api-boilerplate-go/types"
)

type RefreshToken struct {
	ID           int
	UserID       int
	RefreshToken string
	ExpiresAt    time.Time
	Revoked      bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type AuthService struct {
	store types.UserStore
}

func (a *AuthService) CheckUserByEmail(email string) (bool, error) {
	user, err := a.store.GetUserByEmail(email)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return false, nil
		}
		return false, err
	}
	return user != nil, nil
}
