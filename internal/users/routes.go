package users

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ren-zi-fa/rest-api-boilerplate-go/config"
	"github.com/ren-zi-fa/rest-api-boilerplate-go/internal/auth/middleware"
	"github.com/ren-zi-fa/rest-api-boilerplate-go/types"
	"github.com/ren-zi-fa/rest-api-boilerplate-go/utils"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) RegisterRoute(router chi.Router) {
	router.With(middleware.NewAuthMiddleware(config.Envs.JWTSecret), middleware.RoleMiddleware("admin")).Get("/users", h.handleGetUsers)
}
func (h *Handler) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.store.GetUsers()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}
	utils.WriteJSON(w, http.StatusOK, users)
}
