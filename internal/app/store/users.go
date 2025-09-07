package store

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/aseptimu/GophKeeper/internal/app/models"
)

var ErrUserNotFound = errors.New("user not found")

const saveUserQuery = `
	INSERT INTO users (login, password_hash)
	VALUES ($1, $2)
	RETURNING id
`

func (db *DBStore) SaveUser(ctx context.Context, user *models.User) error {
	err := db.pool.QueryRow(ctx, saveUserQuery, user.Login, user.HashedPassword).Scan(&user.ID)
	if err != nil {
		slog.Error("Failed to save user", "err", err)
		return err
	}
	slog.Info("Saved user", "login", user.Login, "id", user.ID)
	return nil
}

const getUserByLoginQuery = `
	SELECT id, login, password_hash FROM users WHERE login = $1
`

func (db *DBStore) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	var user models.User

	err := db.pool.QueryRow(ctx, getUserByLoginQuery, login).Scan(&user.ID, &user.Login, &user.HashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		slog.Error("Failed to get user by login", "login", login, "err", err)
		return nil, err
	}
	return &user, nil
}

func (db *DBStore) GetUserById(_ context.Context) {}
