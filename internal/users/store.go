package users

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

func (s *Store) GetUsers() ([]*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	users := make([]*types.User, 0)
	for rows.Next() {
		user := new(types.User)
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {

	row := s.db.QueryRow("SELECT id, username, email,password, createdAt FROM users WHERE email = ?", email)

	user := new(types.User)

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Store) CreateUser(user *types.User) (int64, error) {
	
	existingUser, err := s.GetUserByEmail(user.Email)
	if err == nil && existingUser != nil {
		return 0, sql.ErrNoRows 
	}
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	res, err := s.db.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)",
		user.Username, user.Email, user.Password)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}
