package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ren-zi-fa/rest-api-boilerplate-go/internal/auth"
	costumMiddleware "github.com/ren-zi-fa/rest-api-boilerplate-go/internal/auth/middleware"
	"github.com/ren-zi-fa/rest-api-boilerplate-go/internal/posts"
	"github.com/ren-zi-fa/rest-api-boilerplate-go/internal/users"
	"github.com/ren-zi-fa/rest-api-boilerplate-go/types"
)

var m types.Middleware = costumMiddleware.MiddlewareImpl{}

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
	router.Use(m.RateLimitMiddleware)
	router.Use(middleware.CleanPath)

	router.Route("/api/auth", func(r chi.Router) {
		usersStore := users.NewStore(s.db)
		authStore := auth.NewStore(s.db)
		authHandler := auth.NewHandler(usersStore, authStore)
		authHandler.RegisterRoute(r)
	})
	router.Route("/api/v1", func(r chi.Router) {

		usersStore := users.NewStore(s.db)
		usersHandler := users.NewHandler(usersStore)
		usersHandler.RegisterRoute(r)

		postStore := posts.NewStore(s.db)
		postsHandler := posts.NewHandler(postStore)
		postsHandler.RegisterRoute(r)

	})

	log.Println("Server running on", s.addr)

	return http.ListenAndServe(s.addr, router)
}
