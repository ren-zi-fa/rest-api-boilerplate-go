package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ren-zi-fa/rest-api-boilerplate-go/config"
	"github.com/ren-zi-fa/rest-api-boilerplate-go/internal/auth/jwt"
	"github.com/ren-zi-fa/rest-api-boilerplate-go/internal/auth/middleware"

	"github.com/ren-zi-fa/rest-api-boilerplate-go/types"
	"github.com/ren-zi-fa/rest-api-boilerplate-go/utils"
)

var m types.Middleware = middleware.MiddlewareImpl{}

type Handler struct {
	user types.UserStore
	auth types.AuthStore
}

func NewHandler(user types.UserStore, auth types.AuthStore) *Handler {
	return &Handler{
		user: user,
		auth: auth,
	}
}

func (h *Handler) RegisterRoute(router chi.Router) {
	router.Post("/login", h.handleLoginUser)
	router.Post("/register", h.handleRegisterUser)
	router.With(m.CheckRefreshToken).Post("/refresh", h.handleRefreshToken)
}

func (h *Handler) handleLoginUser(w http.ResponseWriter, r *http.Request) {
	var loginPayload types.LoginUserPayload
	err := utils.ParseJSON(r, &loginPayload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid request payload"))
		return
	}
	user, err := h.user.GetUserByEmail(loginPayload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user not registered yet"))
		return
	}

	if err := ComparePasswordBcrypt(user.Password, loginPayload.Password); err != nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid email or password"))
		return
	}
	tokenID := uuid.NewString()
	refreshToken, err := jwt.GenerateRefreshToken(uint(user.ID), tokenID, config.Envs.JWTSecret, config.Envs.REFRESH_TOKEN_EXPIRE_DURATION)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed generate refresh token"))
	}
	accessToken, err := jwt.GenerateAccessToken(uint(user.ID), user.Role, config.Envs.JWTSecret, config.Envs.ACCESS_TOKEN_EXPIRE_DURATION)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed generate access token"))
	}
	duration := config.Envs.REFRESH_TOKEN_EXPIRE_DURATION
	expiresAt := time.Now().Add(duration)

	tokenRef := &types.RefreshTokenDB{
		UserID:       user.ID,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}
	resToken, err := h.auth.RefreshTokenStore(tokenRef)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Authorization", "Bearer "+accessToken)

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    resToken.RefreshToken,
		Expires:  resToken.ExpiresAt,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/auth/refresh",
	})
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"access_token": accessToken,
	})
}

func (h *Handler) handleRegisterUser(w http.ResponseWriter, r *http.Request) {
	var userPayload types.RegisterUserPayload
	err := utils.ParseJSON(r, &userPayload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := utils.Validate.Struct(userPayload); err != nil {
		fieldErrors := utils.FormatValidationError(err)
		utils.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"error":  "invalid payload",
			"fields": fieldErrors,
		})
		return
	}

	hashedPassword, err := HashPasswordBcrypt(userPayload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to hash password: %w", err))
		return
	}
	userPayload.Password = hashedPassword

	newUser := &types.User{
		Email:    userPayload.Email,
		Username: userPayload.Username,
		Password: userPayload.Password,
	}
	id, err := h.user.CreateUser(newUser)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("this account has registered"))
		return
	}
	utils.WriteJSON(w, http.StatusCreated, map[string]any{"id": id})

}

func (h *Handler) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	
}
