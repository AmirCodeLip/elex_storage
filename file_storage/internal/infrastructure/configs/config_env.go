package configs

import (
	"elex_storage/file_storage/internal/domain"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// Load configurations from .env file.
func NewConfigEnv() *domain.ConfigEnv {
	config := domain.ConfigEnv{}
	pwd, _ := os.Getwd()
	var exeDir = filepath.Dir(pwd)
	configDir := filepath.Join(exeDir, ".env")
	errEnv := godotenv.Load(configDir)
	if errEnv != nil {
		err := fmt.Errorf("No .env file found")
		log.Fatal(err)
	}
	/// Set config properties
	config.DriveDisk = os.Getenv("DRIVE_DISK")
	config.DriveName = os.Getenv("DRIVE_NAME")
	config.MigrationsDir = os.Getenv("MIGRATIONS_DIR")

	return &config
}

func TestConfigEnv() *domain.ConfigEnv {
	config := domain.ConfigEnv{}
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_DATABASE", "file_storage")
	os.Setenv("DB_USERNAME", "admin")
	os.Setenv("DB_PASSWORD", "admin123")
	os.Setenv("DB_SCHEMA", "public")
	os.Setenv("MIGRATIONS_DIR", "..\\migrations")

	// Set config properties
	config.DriveDisk = os.Getenv("DRIVE_DISK")
	config.DriveName = os.Getenv("DRIVE_NAME")
	config.MigrationsDir = os.Getenv("MIGRATIONS_DIR")
	return &config
}
