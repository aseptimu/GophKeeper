package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserIDKey contextKey = "userID"

func Auth(hmacKey []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			slog.Info("Auth middleware", "path", r.URL.Path, "method", r.Method)

			if r.URL.Path == "/auth/register" || r.URL.Path == "/auth/login" {
				next.ServeHTTP(w, r)
				return
			}

			tokenString := extractToken(r)
			if len(tokenString) > 20 {
				slog.Info("Extracted token", "token", tokenString[:20])
			} else {
				slog.Info("Extracted token", "token", tokenString)
			}
			if tokenString == "" {
				http.Error(w, "Authorization token required", http.StatusUnauthorized)
				return
			}

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return hmacKey, nil
			})

			if err != nil {
				slog.Error("JWT parse error", "error", err)
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			if !token.Valid {
				slog.Error("JWT token invalid")
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				slog.Error("Invalid token claims")
				http.Error(w, "Invalid token claims", http.StatusUnauthorized)
				return
			}

			userID, ok := claims["user_id"].(string)
			if !ok {
				slog.Error("Invalid user ID in token", "user_id", claims["user_id"])
				http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
				return
			}

			slog.Info("JWT validation successful", "user_id", userID)
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func extractToken(r *http.Request) string {
	if cookie, err := r.Cookie("token"); err == nil {
		return cookie.Value
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1]
		}
	}

	return ""
}
