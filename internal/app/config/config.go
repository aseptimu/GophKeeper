package config

import (
	"errors"
	"flag"
	"github.com/aseptimu/GophKeeper/internal/app/utils"
	"github.com/caarlos0/env/v11"
	_ "github.com/joho/godotenv/autoload"
	"log/slog"
	"time"
)

var (
	ErrDatabaseDSNNotSpecified = errors.New("database DSN is not specified")
	ErrJWTKeyTooWeak           = errors.New("JWT_KEY is too short; need >= 32 bytes")
)

const (
	DBTimeout      = 5 * time.Second
	minHMACKeySize = 32
)

type AppConfig struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"127.0.0.1:8087"`
	DSN           string `env:"DATABASE_DSN"`
	JWTKey        string `env:"JWT_KEY"`
	MigrationsDir string `env:"MIGRATIONS_DIR"`
	FilesDir      string `env:"FILES_DIR" envDefault:"./files"`
}

func NewAppConfig() (*AppConfig, error) {
	config := &AppConfig{}

	err := env.Parse(config)
	if err != nil {
		return nil, err
	}

	flag.StringVar(&config.ServerAddress, "a", config.ServerAddress, "Server address")
	flag.StringVar(&config.ServerAddress, "addr", config.ServerAddress, "Server address (alias)")
	flag.StringVar(&config.DSN, "d", config.DSN, "Database DSN")
	flag.StringVar(&config.DSN, "dsn", config.DSN, "Database DSN (alias)")
	flag.StringVar(&config.MigrationsDir, "m", config.MigrationsDir, "Migrations directory")
	flag.StringVar(&config.MigrationsDir, "migration", config.MigrationsDir, "Migrations directory (alias)")
	flag.Parse()

	if config.DSN == "" {
		return nil, ErrDatabaseDSNNotSpecified
	}
	if config.JWTKey == "" {
		config.JWTKey = utils.GenerateRandomSecretKey()
	}
	if len(config.JWTKey) < minHMACKeySize {
		return nil, ErrJWTKeyTooWeak
	}
	if config.MigrationsDir == "" {
		slog.Warn("Migrations directory not specified")
	}

	return config, nil
}
