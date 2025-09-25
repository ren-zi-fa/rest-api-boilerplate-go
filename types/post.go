package types

import "time"

type Post struct {
	ID        int       `json:"id"`
	Author    string    `json:"author"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreatePostPayload struct {
	Author  string `json:"author" validate:"required,min=3"`
	Title   string `json:"title" validate:"required,min=5"`
	Content string `json:"content" validate:"required"`
}

type PostStore interface {
	GetPosts() ([]*Post, error)
	GetPostByID(id int) (*Post, error)
	DeletePostByID(id int) (int64, error)
	CreatePost(CreatePostPayload) (int64, error)
}
