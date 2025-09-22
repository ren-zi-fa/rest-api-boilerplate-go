package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ren-zi-fa/rest-api-boilerplate-go/service/posts"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewServerAPI(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := chi.NewRouter()
	router.Route("/api/v1", func(r chi.Router) {

		postStore := posts.NewStore(s.db)
		postsHandler := posts.NewHandler(postStore)
		postsHandler.RegisterRoute(r)
	})

	log.Println("Server running on", s.addr)

	return http.ListenAndServe(s.addr, router)
}
