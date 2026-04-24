package configs

import (
	"elex_storage/identity_service/internal/domain"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
)

// Load configurations from .env file.
func NewConfigEnv() (*domain.ConfigEnv, error) {
	config := domain.ConfigEnv{}
	pwd, _ := os.Getwd()
	var exeDir = filepath.Dir(pwd)
	configDir := filepath.Join(exeDir, ".env")
	errEnv := godotenv.Load(configDir)
	if errEnv != nil {
		err := fmt.Errorf("No .env file found")
		// return nil, err
		fmt.Println(err)
	}
	// Set AccessTokenDuration
	val1 := os.Getenv("ACCESS_TOKEN_DURATION")
	accessTokenDuration, err1 := time.ParseDuration(val1)
	if err1 != nil {
		return nil, err1
	}
	val2 := os.Getenv("REFRESH_TOKEN_DURATION")
	refreshTokenDuration, err2 := time.ParseDuration(val2)
	if err2 != nil {
		return nil, err2
	}
	/// Set config properties
	config.RefreshTokenDuration = refreshTokenDuration
	config.AccessTokenDuration = accessTokenDuration
	config.MigrationsDir = os.Getenv("MIGRATIONS_DIR")
	config.GrpcPort = os.Getenv("GRPC_PORT")

	return &config, nil
}

func TestConfigEnv() *domain.ConfigEnv {
	config := domain.ConfigEnv{}
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_DATABASE", "identity_service")
	os.Setenv("DB_USERNAME", "admin")
	os.Setenv("DB_PASSWORD", "admin123")
	os.Setenv("DB_SCHEMA", "public")
	os.Setenv("MIGRATIONS_DIR", "..\\migrations")

	/// Set config properties
	config.MigrationsDir = os.Getenv("MIGRATIONS_DIR")
	return &config
}
