package auth

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ren-zi-fa/rest-api-boilerplate-go/types"
	"github.com/ren-zi-fa/rest-api-boilerplate-go/utils"
)

type Handler struct {
	store types.UserStore
	auth  types.Auth
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{
		store: store,
		auth:  &AuthService{store: store},
	}
}

func (h *Handler) RegisterRoute(router chi.Router) {
	router.Post("/login", h.handleLoginUser)
	router.Post("/register", h.handleRegisterUser)
}

func (h *Handler) handleLoginUser(w http.ResponseWriter, r *http.Request) {
	var loginPayload types.LoginUserPayload
	err := utils.ParseJSON(r, &loginPayload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid request payload"))
		return
	}
	user, err := h.store.GetUserByEmail(loginPayload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user not registered yet"))
		return
	}

	fmt.Println("from user", loginPayload.Password)
	if err := utils.ComparePasswordBcrypt(user.Password, loginPayload.Password); err != nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid email or password"))
		return
	}

	//TODO: generate token

}

func (h *Handler) handleRegisterUser(w http.ResponseWriter, r *http.Request) {
	var user types.RegisterUserPayload
	err := utils.ParseJSON(r, &user)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := utils.Validate.Struct(user); err != nil {
		fieldErrors := utils.FormatValidationError(err)
		utils.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"error":  "invalid payload",
			"fields": fieldErrors,
		})
		return
	}

	exists, err := h.auth.CheckUserByEmail(user.Email)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to check email: %w", err))
		return
	}

	if exists {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("email already registered"))
		return
	}

	hashedPassword, err := utils.HashPasswordBcrypt(user.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to hash password: %w", err))
		return
	}
	user.Password = hashedPassword

	newUser := &types.User{
		Email:    user.Email,
		Username: user.Username,
		Password: user.Password,
	}
	id, err := h.store.CreateUser(newUser)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, map[string]any{"id": id})

}
