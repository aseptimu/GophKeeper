package services

import (
	"context"
	"errors"
	"testing"

	"github.com/aseptimu/GophKeeper/internal/app/config"
	"github.com/aseptimu/GophKeeper/internal/app/models"
	"github.com/aseptimu/GophKeeper/internal/app/store"
)

type mockUserStore struct {
	users map[string]*models.User
}

func newMockUserStore() *mockUserStore {
	return &mockUserStore{
		users: make(map[string]*models.User),
	}
}

func (m *mockUserStore) GetUserById(_ context.Context) {
	// Not implemented for this test
}

func (m *mockUserStore) GetUserByLogin(_ context.Context, login string) (*models.User, error) {
	user, exists := m.users[login]
	if !exists {
		return nil, store.ErrUserNotFound
	}
	return user, nil
}

func (m *mockUserStore) SaveUser(_ context.Context, user *models.User) error {
	m.users[user.Login] = user
	return nil
}

func TestAuthService_RegisterUser(t *testing.T) {
	cfg := &config.AppConfig{
		JWTKey: "test-jwt-key-that-is-long-enough-for-hmac",
	}
	mockStore := newMockUserStore()
	authService := NewAuthService(cfg, mockStore)

	tests := []struct {
		name        string
		login       string
		password    string
		expectError bool
		errorType   error
	}{
		{
			name:        "successful registration",
			login:       "testuser",
			password:    "password123",
			expectError: false,
		},
		{
			name:        "empty login",
			login:       "",
			password:    "password123",
			expectError: true,
			errorType:   ErrEmptyCredentials,
		},
		{
			name:        "empty password",
			login:       "testuser",
			password:    "",
			expectError: true,
			errorType:   ErrEmptyCredentials,
		},
		{
			name:        "user already exists",
			login:       "existinguser",
			password:    "password123",
			expectError: true,
			errorType:   ErrUserAlreadyExists,
		},
	}

	// Pre-register a user for the "user already exists" test
	existingUser := &models.User{
		ID:             "existing-id",
		Login:          "existinguser",
		HashedPassword: "hashed-password",
	}
	mockStore.SaveUser(context.Background(), existingUser)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := authService.RegisterUser(context.Background(), tt.login, tt.password)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if tt.errorType != nil && !errors.Is(err, tt.errorType) {
					t.Errorf("expected error type %v, got %v", tt.errorType, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if token == "" {
					t.Errorf("expected token but got empty string")
				}
			}
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	cfg := &config.AppConfig{
		JWTKey: "test-jwt-key-that-is-long-enough-for-hmac",
	}
	mockStore := newMockUserStore()
	authService := NewAuthService(cfg, mockStore)

	// Register a test user
	_, err := authService.RegisterUser(context.Background(), "testuser", "password123")
	if err != nil {
		t.Fatalf("failed to register test user: %v", err)
	}

	tests := []struct {
		name        string
		login       string
		password    string
		expectError bool
		errorType   error
	}{
		{
			name:        "successful login",
			login:       "testuser",
			password:    "password123",
			expectError: false,
		},
		{
			name:        "empty login",
			login:       "",
			password:    "password123",
			expectError: true,
			errorType:   ErrEmptyCredentials,
		},
		{
			name:        "empty password",
			login:       "testuser",
			password:    "",
			expectError: true,
			errorType:   ErrEmptyCredentials,
		},
		{
			name:        "user not found",
			login:       "nonexistent",
			password:    "password123",
			expectError: true,
			errorType:   store.ErrUserNotFound,
		},
		{
			name:        "invalid password",
			login:       "testuser",
			password:    "wrongpassword",
			expectError: true,
			errorType:   ErrInvalidPassword,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := authService.Login(context.Background(), tt.login, tt.password)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if tt.errorType != nil && !errors.Is(err, tt.errorType) {
					t.Errorf("expected error type %v, got %v", tt.errorType, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if token == "" {
					t.Errorf("expected token but got empty string")
				}
			}
		})
	}
}

func TestAuthService_GenerateToken(t *testing.T) {
	cfg := &config.AppConfig{
		JWTKey: "test-jwt-key-that-is-long-enough-for-hmac",
	}
	mockStore := newMockUserStore()
	authService := NewAuthService(cfg, mockStore)

	token, err := authService.generateToken("test-user-id", cfg.JWTKey)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	if token == "" {
		t.Errorf("expected token but got empty string")
	}

	// Test with different user ID
	token2, err := authService.generateToken("test-user-id-2", cfg.JWTKey)
	if err != nil {
		t.Fatalf("failed to generate second token: %v", err)
	}

	if token == token2 {
		t.Errorf("tokens should be different for different user IDs")
	}
}
