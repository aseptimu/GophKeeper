package config

import (
	"os"
	"testing"

	"github.com/aseptimu/GophKeeper/internal/app/utils"
	"github.com/caarlos0/env/v11"
)

func TestNewAppConfig(t *testing.T) {
	// Save original environment
	originalDSN := os.Getenv("DATABASE_DSN")
	originalJWTKey := os.Getenv("JWT_KEY")
	originalServerAddr := os.Getenv("SERVER_ADDRESS")
	originalMigrationsDir := os.Getenv("MIGRATIONS_DIR")
	originalFilesDir := os.Getenv("FILES_DIR")

	// Clean up after test
	defer func() {
		if originalDSN != "" {
			os.Setenv("DATABASE_DSN", originalDSN)
		} else {
			os.Unsetenv("DATABASE_DSN")
		}
		if originalJWTKey != "" {
			os.Setenv("JWT_KEY", originalJWTKey)
		} else {
			os.Unsetenv("JWT_KEY")
		}
		if originalServerAddr != "" {
			os.Setenv("SERVER_ADDRESS", originalServerAddr)
		} else {
			os.Unsetenv("SERVER_ADDRESS")
		}
		if originalMigrationsDir != "" {
			os.Setenv("MIGRATIONS_DIR", originalMigrationsDir)
		} else {
			os.Unsetenv("MIGRATIONS_DIR")
		}
		if originalFilesDir != "" {
			os.Setenv("FILES_DIR", originalFilesDir)
		} else {
			os.Unsetenv("FILES_DIR")
		}
	}()

	tests := []struct {
		name        string
		envVars     map[string]string
		expectError bool
		errorType   error
	}{
		{
			name: "valid config with all env vars",
			envVars: map[string]string{
				"DATABASE_DSN":   "postgres://user:pass@localhost:5432/db",
				"JWT_KEY":        "test-jwt-key-that-is-long-enough-for-hmac",
				"SERVER_ADDRESS": "0.0.0.0:8080",
				"MIGRATIONS_DIR": "/app/migrations",
				"FILES_DIR":      "/app/files",
			},
			expectError: false,
		},
		{
			name: "valid config with minimal env vars",
			envVars: map[string]string{
				"DATABASE_DSN": "postgres://user:pass@localhost:5432/db",
			},
			expectError: false,
		},
		{
			name: "missing database DSN",
			envVars: map[string]string{
				"JWT_KEY": "test-jwt-key-that-is-long-enough-for-hmac",
			},
			expectError: true,
			errorType:   ErrDatabaseDSNNotSpecified,
		},
		{
			name: "JWT key too short",
			envVars: map[string]string{
				"DATABASE_DSN": "postgres://user:pass@localhost:5432/db",
				"JWT_KEY":      "short",
			},
			expectError: true,
			errorType:   ErrJWTKeyTooWeak,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment
			os.Unsetenv("DATABASE_DSN")
			os.Unsetenv("JWT_KEY")
			os.Unsetenv("SERVER_ADDRESS")
			os.Unsetenv("MIGRATIONS_DIR")
			os.Unsetenv("FILES_DIR")

			// Set test environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			config, err := newAppConfigForTest()

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if tt.errorType != nil && err != tt.errorType {
					t.Errorf("expected error type %v, got %v", tt.errorType, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if config == nil {
					t.Fatal("expected config but got nil")
				}

				// Check that environment variables are properly set
				if tt.envVars["DATABASE_DSN"] != "" && config.DSN != tt.envVars["DATABASE_DSN"] {
					t.Errorf("expected DSN %s, got %s", tt.envVars["DATABASE_DSN"], config.DSN)
				}
				if tt.envVars["JWT_KEY"] != "" && config.JWTKey != tt.envVars["JWT_KEY"] {
					t.Errorf("expected JWTKey %s, got %s", tt.envVars["JWT_KEY"], config.JWTKey)
				}
				if tt.envVars["SERVER_ADDRESS"] != "" && config.ServerAddress != tt.envVars["SERVER_ADDRESS"] {
					t.Errorf("expected ServerAddress %s, got %s", tt.envVars["SERVER_ADDRESS"], config.ServerAddress)
				}
				if tt.envVars["MIGRATIONS_DIR"] != "" && config.MigrationsDir != tt.envVars["MIGRATIONS_DIR"] {
					t.Errorf("expected MigrationsDir %s, got %s", tt.envVars["MIGRATIONS_DIR"], config.MigrationsDir)
				}
				if tt.envVars["FILES_DIR"] != "" && config.FilesDir != tt.envVars["FILES_DIR"] {
					t.Errorf("expected FilesDir %s, got %s", tt.envVars["FILES_DIR"], config.FilesDir)
				}

				// Check defaults
				if tt.envVars["SERVER_ADDRESS"] == "" && config.ServerAddress != "127.0.0.1:8087" {
					t.Errorf("expected default ServerAddress 127.0.0.1:8087, got %s", config.ServerAddress)
				}
				if tt.envVars["FILES_DIR"] == "" && config.FilesDir != "./files" {
					t.Errorf("expected default FilesDir ./files, got %s", config.FilesDir)
				}
			}
		})
	}
}

func TestAppConfig_Defaults(t *testing.T) {
	// Save original environment
	originalDSN := os.Getenv("DATABASE_DSN")
	originalJWTKey := os.Getenv("JWT_KEY")
	originalServerAddr := os.Getenv("SERVER_ADDRESS")
	originalMigrationsDir := os.Getenv("MIGRATIONS_DIR")
	originalFilesDir := os.Getenv("FILES_DIR")

	// Clean up after test
	defer func() {
		if originalDSN != "" {
			os.Setenv("DATABASE_DSN", originalDSN)
		} else {
			os.Unsetenv("DATABASE_DSN")
		}
		if originalJWTKey != "" {
			os.Setenv("JWT_KEY", originalJWTKey)
		} else {
			os.Unsetenv("JWT_KEY")
		}
		if originalServerAddr != "" {
			os.Setenv("SERVER_ADDRESS", originalServerAddr)
		} else {
			os.Unsetenv("SERVER_ADDRESS")
		}
		if originalMigrationsDir != "" {
			os.Setenv("MIGRATIONS_DIR", originalMigrationsDir)
		} else {
			os.Unsetenv("MIGRATIONS_DIR")
		}
		if originalFilesDir != "" {
			os.Setenv("FILES_DIR", originalFilesDir)
		} else {
			os.Unsetenv("FILES_DIR")
		}
	}()

	// Clear environment
	os.Unsetenv("DATABASE_DSN")
	os.Unsetenv("JWT_KEY")
	os.Unsetenv("SERVER_ADDRESS")
	os.Unsetenv("MIGRATIONS_DIR")
	os.Unsetenv("FILES_DIR")

	// Set only required DSN
	os.Setenv("DATABASE_DSN", "postgres://user:pass@localhost:5432/db")

	config, err := newAppConfigForTest()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check defaults
	if config.ServerAddress != "127.0.0.1:8087" {
		t.Errorf("expected default ServerAddress 127.0.0.1:8087, got %s", config.ServerAddress)
	}
	if config.FilesDir != "./files" {
		t.Errorf("expected default FilesDir ./files, got %s", config.FilesDir)
	}
	if config.MigrationsDir != "" {
		t.Errorf("expected empty MigrationsDir, got %s", config.MigrationsDir)
	}
	if config.JWTKey == "" {
		t.Error("expected JWTKey to be generated, got empty string")
	}
	if len(config.JWTKey) < minHMACKeySize {
		t.Errorf("expected JWTKey length >= %d, got %d", minHMACKeySize, len(config.JWTKey))
	}
}

// newAppConfigForTest creates a new AppConfig for testing without command line flags
func newAppConfigForTest() (*AppConfig, error) {
	config := &AppConfig{}

	// Parse environment variables only
	if err := env.Parse(config); err != nil {
		return nil, err
	}

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
		// No warning in tests
	}

	return config, nil
}
