package auth

import (
	"database/sql"

	"github.com/ren-zi-fa/rest-api-boilerplate-go/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) RefreshTokenStore(token *types.RefreshTokenDB) (*types.RefreshTokenDB, error) {
	query := `
        INSERT INTO refresh_tokens (user_id, refresh_token, expires_at, revoked)
        VALUES (?, ?, ?, ?)
    `

	result, err := s.db.Exec(query, token.UserID, token.RefreshToken, token.ExpiresAt, false)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	token.ID = int(id)
	return token, nil
}
