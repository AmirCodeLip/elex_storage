package database

import (
	"elex_storage/pkg/logger"
	"elex_storage/pkg/shared_kernel/models"
	"fmt"

	"elex_storage/pkg/auto_migration"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
)

func NewDatabase(config *models.ConfigEnv2, logger logger.Logger) *sqlx.DB {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable&search_path=%s",
		config.Database.User,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.Name,
		config.Database.Schema)
	db, err := sqlx.Open("pgx", connStr)
	if err != nil {
		logger.Error(err.Error())
	}
	return db
}

func ApplyMigration(db *sqlx.DB, cfg *models.ConfigEnv2, logger logger.Logger) error {
	autoMigrateManager, err := auto_migration.NewAutoMigrateManager(db, cfg.MigrationsDir)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	// Apply Migrations
	migrationsErr := autoMigrateManager.AutoMigrate()
	if migrationsErr != nil {
		fmt.Println(migrationsErr.Error())
		logger.Error(migrationsErr.Error())
		return err
	}
	return nil
}
