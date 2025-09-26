package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ren-zi-fa/rest-api-boilerplate-go/config"
	"github.com/stretchr/testify/assert"
)

func TestGenerateAndParseRefreshToken(t *testing.T) {
	userID := uint(1)
	secretKey := config.Envs.JWTSecret
	duration := config.Envs.REFRESH_TOKEN_EXPIRE_DURATION
	tokenID := uuid.NewString()

	token, err := GenerateRefreshToken(userID, tokenID, secretKey, duration)

	assert.NoError(t, err, "should not error when generating refresh token")
	assert.NotEmpty(t, token, "refresh token should not be empty")

	claims, err := ParseRefreshToken(token, secretKey)

	assert.NoError(t, err, "valid token should be parsed successfully")
	assert.NotNil(t, claims, "claims should not be nil")
	assert.Equal(t, userID, claims.UserID, "UserID in claims should match")
	assert.Equal(t, tokenID, claims.TokenID, "TokenID in claims should match")
	assert.Equal(t, "refresh_token", claims.Subject, "Subject should be 'refresh_token'")

	expectedExpiry := time.Now().Add(duration)
	actualExpiry := claims.ExpiresAt.Time
	assert.WithinDuration(t,
		expectedExpiry,
		actualExpiry,
		5*time.Second,
		"ExpiredAt should match duration")

}

func TestGenerateAndParseAcessToken(t *testing.T) {
	userID := uint(2)
	secretKey := config.Envs.JWTSecret
	duration := config.Envs.ACCESS_TOKEN_EXPIRE_DURATION

	token, err := GenerateAccessToken(userID, secretKey, duration)
	assert.NoError(t, err, "should not error when generating access token")
	assert.NotEmpty(t, token, "access token should not be empty")

	claims, err := ParseAccessToken(token, secretKey)

	assert.NoError(t, err, "valid token should be parsed successfully")
	assert.NotNil(t, claims, "claims should not be nil")
	assert.Equal(t, userID, claims.UserID, "UserID in claims should match")
	assert.Equal(t, "access_token", claims.Subject, "Subject should be 'access_token'")

	expectedExpiry := time.Now().Add(duration)
	actualExpiry := claims.ExpiresAt.Time
	assert.WithinDuration(t,
		expectedExpiry,
		actualExpiry,
		5*time.Second,
		"ExpiredAt should match duration")
}
