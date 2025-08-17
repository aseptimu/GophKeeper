package middleware

import (
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := jwt.New(jwt.SigningMethodHS256)
		s, _ := t.SignedString([]byte("test"))
		w.Header().Add("Content-Type", "application/json")
		http.SetCookie(w, &http.Cookie{
			Name:  "access_token",
			Value: s,
		})
		next.ServeHTTP(w, r)
	})
}
