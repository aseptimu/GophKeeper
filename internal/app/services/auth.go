package services

import (
	"context"
	"errors"
	"github.com/aseptimu/GophKeeper/internal/app/config"
	"github.com/aseptimu/GophKeeper/internal/app/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

const tokenExpirationTime = time.Hour * 72

type UserStore interface {
	GetUserById(ctx context.Context)
	GetUserByLogin(ctx context.Context, login string) (*models.User, error)
	SaveUser(ctx context.Context, user *models.User) error
}

type AuthService struct {
	store UserStore
	cfg   *config.AppConfig
}

func NewAuthService(cfg *config.AppConfig, store UserStore) *AuthService {
	return &AuthService{
		store: store,
		cfg:   cfg,
	}
}

func (s *AuthService) RegisterUser(ctx context.Context, login, password string) (string, error) {
	if login == "" || password == "" {
		return "", ErrEmptyCredentials
	}

	existingUser, err := s.store.GetUserByLogin(ctx, login)
	if err != nil && !errors.Is(err, ErrUserNotFound) {
		return "", err
	}
	if existingUser != nil {
		return "", ErrUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := &models.User{
		Login:          login,
		HashedPassword: string(hashedPassword),
	}
	err = s.store.SaveUser(ctx, user)
	if err != nil {
		return "", err
	}

	token, err := s.generateToken(user.ID, s.cfg.JWTKey)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *AuthService) Login(ctx context.Context, login, password string) (string, error) {
	if login == "" || password == "" {
		return "", ErrEmptyCredentials
	}

	user, err := s.store.GetUserByLogin(ctx, login)
	if err != nil {
		return "", err
	}

	token, err := s.generateToken(user.ID, s.cfg.JWTKey)
	if err != nil {
		slog.Error("Failed to generate jwt token", "error", err.Error())
		return "", err
	}
	return token, nil
}

func (s *AuthService) generateToken(userID string, secretKey string) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     jwt.NewNumericDate(time.Now().Add(tokenExpirationTime)),
	})
	return t.SignedString([]byte(secretKey))
}
