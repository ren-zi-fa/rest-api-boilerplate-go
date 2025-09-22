package posts

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

func (s *Store) GetPosts() ([]*types.Post, error) {
	rows, err := s.db.Query("SELECT * FROM post")
	if err != nil {
		return nil, err
	}

	posts := make([]*types.Post, 0)
	for rows.Next() {
		post, err := scanRowsIntoPost(rows)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (s *Store) GetPostByID(postID int) (*types.Post, error) {

	row := s.db.QueryRow("SELECT id, title, author, content, createdAt FROM post WHERE id = ?", postID)

	post := new(types.Post)

	err := row.Scan(
		&post.ID,
		&post.Title,
		&post.Author,
		&post.Content,
		&post.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (s *Store) DeletePostByID(id int) (int64, error) {

	res, err := s.db.Exec("DELETE FROM post WHERE id = ?", id)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func (s *Store) CreatePost(payload types.CreatePostPayload) (int64, error) {

	result, err := s.db.Exec("INSERT INTO post (title, author, content) VALUES (?, ?, ?)",
		payload.Title, payload.Author, payload.Content)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {

		return 0, err
	}
	return id, nil
}

// scan row
func scanRowsIntoPost(rows *sql.Rows) (*types.Post, error) {
	post := new(types.Post)

	err := rows.Scan(
		&post.ID,
		&post.Title,
		&post.Author,
		&post.Content,
		&post.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return post, nil
}
