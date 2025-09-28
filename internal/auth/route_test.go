package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/ren-zi-fa/rest-api-boilerplate-go/config"
	jwtCostum "github.com/ren-zi-fa/rest-api-boilerplate-go/internal/auth/jwt"
	"github.com/ren-zi-fa/rest-api-boilerplate-go/internal/auth/middleware"
	"github.com/ren-zi-fa/rest-api-boilerplate-go/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUserRepo struct{ mock.Mock }

func (m *mockUserRepo) GetUserByEmail(email string) (*types.User, error) {
	args := m.Called(email)
	return args.Get(0).(*types.User), args.Error(1)
}
func (m *mockUserRepo) CreateUser(user *types.User) (int64, error) {
	args := m.Called(user)
	return args.Get(0).(int64), args.Error(1)
}
func (m *mockUserRepo) GetUsers() ([]*types.User, error) {
	args := m.Called()
	return args.Get(0).([]*types.User), args.Error(1)
}

type mockAuthRepo struct{ mock.Mock }

func (m *mockAuthRepo) RefreshTokenStore(token *types.RefreshTokenDB) (*types.RefreshTokenDB, error) {
	args := m.Called(token)
	return args.Get(0).(*types.RefreshTokenDB), args.Error(1)
}

func (m *mockAuthRepo) GetRefreshTokenByTokenID(tokenID string) (*types.RefreshTokenDB, error) {
	args := m.Called(tokenID)
	return args.Get(0).(*types.RefreshTokenDB), args.Error(1)
}
func (m *mockAuthRepo) RevokeRefreshToken(userID int, tokenID string) error {
	args := m.Called(userID, tokenID)
	return args.Error(0)
}

func TestHandleLoginUser_Success(t *testing.T) {

	userRepo := new(mockUserRepo)
	authRepo := new(mockAuthRepo)

	h := &Handler{
		user: userRepo,
		auth: authRepo,
	}

	hashedPwd, _ := HashPasswordBcrypt("password123")

	user := &types.User{
		ID:       1,
		Email:    "test@example.com",
		Password: hashedPwd,
	}

	userRepo.On("GetUserByEmail", "test@example.com").Return(user, nil)

	dummyRefresh := &types.RefreshTokenDB{
		UserID:       user.ID,
		RefreshToken: uuid.NewString(),
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
	}
	authRepo.On("RefreshTokenStore", mock.Anything).Return(dummyRefresh, nil)
	authRepo.On("RevokeRefreshToken", user.ID, mock.AnythingOfType("string")).Return(nil)

	payload := types.LoginUserPayload{
		Email:    "test@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.handleLoginUser(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.NotEmpty(t, res.Header.Get("Authorization"))
	assert.NotEmpty(t, res.Cookies())

	var respBody map[string]string
	_ = json.NewDecoder(res.Body).Decode(&respBody)
	assert.NotEmpty(t, respBody["access_token"])
}

func TestHandleRefreshToken_Success(t *testing.T) {
	userRepo := new(mockUserRepo)
	authRepo := new(mockAuthRepo)

	h := &Handler{
		user: userRepo,
		auth: authRepo,
	}

	duration := config.Envs.REFRESH_TOKEN_EXPIRE_DURATION

	refreshClaims := &jwtCostum.RefreshClaims{
		UserID:  1,
		TokenID: uuid.NewString(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "refresh_token",
		},
	}


	newTokenStr := "new_refresh_token"
	newToken := &types.RefreshTokenDB{
		UserID:       1,
		RefreshToken: newTokenStr,
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
	}


	authRepo.On("RefreshTokenStore", mock.AnythingOfType("*types.RefreshTokenDB")).Return(newToken, nil).Twice()


	authRepo.On("RevokeRefreshToken", int(refreshClaims.UserID), refreshClaims.TokenID).Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)

	ctx := req.Context()
	ctx = context.WithValue(ctx, middleware.RefreshToken, refreshClaims)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	h.handleRefreshToken(w, req)

	res := w.Result()
	defer res.Body.Close()


	assert.Equal(t, http.StatusOK, res.StatusCode)

	cookies := res.Cookies()
	assert.NotEmpty(t, cookies)

	found := false
	for _, c := range cookies {
		if c.Name == "refresh_token" && c.Value == newTokenStr {
			found = true
		}
	}
	assert.True(t, found, "refresh_token cookie should be set")

	authRepo.AssertExpectations(t)
}
