package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/aseptimu/GophKeeper/internal/app/services"
	"github.com/aseptimu/GophKeeper/internal/app/store"
	"github.com/go-chi/chi/v5"
)

type AuthUser interface {
	Login(ctx context.Context, login, password string) (string, error)
	RegisterUser(ctx context.Context, login, password string) (string, error)
}

type AuthHandler struct {
	srv AuthUser
}

func NewAuthHandler(srv AuthUser) *AuthHandler {
	return &AuthHandler{
		srv: srv,
	}
}

func (h *AuthHandler) RegisterRoutes(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", h.Register)
		r.Post("/login", h.Login)
	})
}

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var user User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	token, err := h.srv.Login(r.Context(), user.Login, user.Password)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrEmptyCredentials):
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		case errors.Is(err, store.ErrUserNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		case errors.Is(err, services.ErrInvalidPassword):
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "token",
		Value: token,
	})
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	token, err := h.srv.RegisterUser(r.Context(), user.Login, user.Password)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrEmptyCredentials):
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		case errors.Is(err, services.ErrUserAlreadyExists):
			http.Error(w, err.Error(), http.StatusConflict)
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "token",
		Value: token,
	})
}
