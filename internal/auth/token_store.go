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

func (s *Store) GetRefreshTokenByTokenID(tokenID string) (*types.RefreshTokenDB, error) {
	query := `
        SELECT id, user_id, refresh_token, expires_at, revoked
        FROM refresh_tokens
        WHERE id = ?
    `
	row := s.db.QueryRow(query, tokenID)

	var token types.RefreshTokenDB
	err := row.Scan(
		&token.ID,
		&token.UserID,
		&token.RefreshToken,
		&token.ExpiresAt,
		&token.Revoked,
	)
	if err != nil {
		return nil, err
	}

	return &token, nil
}
