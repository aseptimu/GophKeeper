package middleware

import (
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

func Auth(hmacKey []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"userId": "1",
			})
			s, _ := t.SignedString(hmacKey)
			w.Header().Add("Content-Type", "application/json")
			http.SetCookie(w, &http.Cookie{
				Name:  "access_token",
				Value: s,
			})
			next.ServeHTTP(w, r)
		})
	}
}
