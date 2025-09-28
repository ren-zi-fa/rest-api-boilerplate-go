package middleware

import (
	"context"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/ren-zi-fa/rest-api-boilerplate-go/config"
	"github.com/ren-zi-fa/rest-api-boilerplate-go/internal/auth/jwt"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

type MiddlewareImpl struct{}
type contextKey string

const (
	ContextUserID contextKey = "userID"
	ContextRole   contextKey = "role"
	RefreshToken  contextKey = "refreshToken"
)

func (m MiddlewareImpl) NewAuthMiddleware(secretKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "missing token", http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := jwt.ParseAccessToken(tokenStr, secretKey)
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), ContextUserID, claims.UserID)
			ctx = context.WithValue(ctx, ContextRole, claims.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (m MiddlewareImpl) RoleMiddleware(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := r.Context().Value(ContextRole).(string)
			if !ok {
				http.Error(w, "role not found", http.StatusForbidden)
				return
			}

			if slices.Contains(allowedRoles, role) {
				next.ServeHTTP(w, r)
				return
			}

			http.Error(w, "forbidden", http.StatusForbidden)
		})
	}
}

// RateLimitMiddleware limits the number of requests from a single IP address.

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	clients         = make(map[string]*client)
	mu              sync.Mutex
	rateLimit       = rate.Limit(10)
	burst           = 10
	cleanUpInterval = time.Minute * 5
)

func init() {
	go cleanupClients()
}

func cleanupClients() {
	for {
		time.Sleep(cleanUpInterval)
		mu.Lock()
		for ip, c := range clients {
			if time.Since(c.lastSeen) > 10*time.Minute {
				delete(clients, ip)
			}
		}
		mu.Unlock()
	}
}

func (m MiddlewareImpl) RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.Header.Get("X-Forwarded-For")
		if ip == "" {
			host, _, err := net.SplitHostPort(r.RemoteAddr)
			if err == nil {
				ip = host
			} else {
				ip = r.RemoteAddr
			}
		}

		mu.Lock()
		c, exists := clients[ip]
		if !exists {
			limiter := rate.NewLimiter(rateLimit, burst)
			c = &client{limiter, time.Now()}
			clients[ip] = c
		}
		c.lastSeen = time.Now()
		mu.Unlock()

		if !c.limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (m MiddlewareImpl) CheckRefreshToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("refresh_token")
		if err != nil {
			http.Error(w, "Refresh token missing", http.StatusUnauthorized)
			return
		}

		claims, err := jwt.ValidateRefreshToken(cookie.Value, config.Envs.JWTSecret)
		if err != nil {
			http.Error(w, "Invalid refresh token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), RefreshToken, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m MiddlewareImpl) LoggerJSON(next http.Handler) http.Handler {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})


	logFile := filepath.Join("log", "log.file")

	
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Info("Gagal membuka file log, menggunakan stdout")
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(rw, r)

		log.WithFields(logrus.Fields{
			"method":       r.Method,
			"uri":          r.RequestURI,
			"status":       rw.Status(),
			"remote_ip":    r.RemoteAddr,
			"user_agent":   r.UserAgent(),
			"duration_ms":  time.Since(start).Milliseconds(),
			"duration_ns":  time.Since(start).Nanoseconds(),
			"request_time": start.Format(time.RFC3339),
		}).Info("request log")
	})
}
