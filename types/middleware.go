package types

import "net/http"

type Middleware interface {
	RateLimitMiddleware(next http.Handler) http.Handler
	RoleMiddleware(allowedRoles ...string) func(http.Handler) http.Handler
	NewAuthMiddleware(secretKey string) func(http.Handler) http.Handler
	CheckRefreshToken(next http.Handler) http.Handler
}
