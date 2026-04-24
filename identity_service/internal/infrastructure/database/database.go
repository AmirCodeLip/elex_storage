package database

import (
	"elex_storage/identity_service/internal/domain"
	"elex_storage/pkg/logger"
	"fmt"
	"os"

	"elex_storage/pkg/auto_migration"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
)

var (
	database string
	password string
	username string
	port     string
	host     string
	schema   string
	db       *sqlx.DB
)

func NewDatabase(config *domain.ConfigEnv, logger logger.Logger) *sqlx.DB {
	database = os.Getenv("DB_DATABASE")
	password = os.Getenv("DB_PASSWORD")
	username = os.Getenv("DB_USERNAME")
	port = os.Getenv("DB_PORT")
	host = os.Getenv("DB_HOST")
	schema = os.Getenv("DB_SCHEMA")

	// Reuse Connection
	if db != nil {
		return db
	}
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", username, password, host, port, database, schema)
	db, err := sqlx.Open("pgx", connStr)
	if err != nil {
		logger.Error(err.Error())
	}
	return db
}

func ApplyMigration(db *sqlx.DB, cfg *domain.ConfigEnv, logger logger.Logger) error {
	autoMigrateManager, err := auto_migration.NewAutoMigrateManager(db, cfg.MigrationsDir)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	// Apply Migrations
	migrationsErr := autoMigrateManager.AutoMigrate()
	if migrationsErr != nil {
		logger.Error(migrationsErr.Error())
		return migrationsErr
	}
	return nil
}
