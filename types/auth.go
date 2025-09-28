package types

import "time"

type RefreshTokenDB struct {
	ID           int
	UserID       int
	RefreshToken string
	ExpiresAt    time.Time
	Revoked      bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type AuthStore interface {
	RefreshTokenStore(token *RefreshTokenDB) (*RefreshTokenDB, error)
	GetRefreshTokenByTokenID(tokenID string) (*RefreshTokenDB, error)
	RevokeRefreshToken(userId int, token string) error
}
