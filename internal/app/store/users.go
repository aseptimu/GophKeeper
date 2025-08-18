package store

import (
	"context"
	"database/sql"
	"errors"
	"github.com/aseptimu/GophKeeper/internal/app/models"
	"github.com/aseptimu/GophKeeper/internal/app/services"
	"log/slog"
)

const saveUserQuery = `
	INSERT INTO users (login, password_hash)
	VALUES ($1, $2)
`

func (db *DBStore) SaveUser(ctx context.Context, user *models.User) error {
	status, err := db.pool.Exec(ctx, saveUserQuery, user.Login, user.HashedPassword)
	if err != nil {
		slog.Error("Failed to save user", "err", err)
		return err
	}
	slog.Info("Saved user", "login", user.Login, "query status", status)
	return nil
}

const getUserByLoginQuery = `
	SELECT login, password_hash FROM users WHERE login = $1
`

func (db *DBStore) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	var user models.User

	err := db.pool.QueryRow(ctx, getUserByLoginQuery, login).Scan(&user.Login, &user.HashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, services.ErrUserNotFound
		}
		slog.Error("Failed to get user by login", "login", login, "err", err)
		return nil, err
	}
	return &user, nil
}

func (db *DBStore) GetUserById(ctx context.Context) {}
