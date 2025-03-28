package server

import (
	"context"
	"net/http"
	"strings"

	"github.com/fmich7/fyle/pkg/auth"
)

// https://golang.org/pkg/context/#WithValue
type CtxUsernameKey struct{}

// AuthMiddleware extracts jwt authorization token from headers and passes it in ctx.
func (s *Server) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		// parse bearer <token> format
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		claims, err := auth.ParseToken(s.jwtSecretKey, tokenString)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), CtxUsernameKey{}, claims.Username)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
