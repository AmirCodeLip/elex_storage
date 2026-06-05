package database

import (
	"context"
	"elex_storage/pkg/logger"
	"elex_storage/pkg/shared_kernel/models"
	"fmt"
	"testing"
	"time"

	"elex_storage/pkg/auto_migration"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func NewDatabase(config *models.ConfigEnv, logger logger.Logger) *sqlx.DB {
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

func NewTestDatabase(t *testing.T, cfg *models.ConfigEnv) *sqlx.DB {
	ctx := context.Background()

	// Start a fresh PostgreSQL container
	pgContainer, err := postgres.Run(
		ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(cfg.Database.Name),
		postgres.WithUsername(cfg.Database.User),
		postgres.WithPassword(cfg.Database.Password),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	testcontainers.CleanupContainer(t, pgContainer) // automatic cleanup
	require.NoError(t, err)

	// Get connection string
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	// Connect with sqlx
	db, err := sqlx.Open("pgx", connStr)
	require.NoError(t, err)

	// Optional: Verify connection
	err = db.Ping()
	require.NoError(t, err)

	return db
}

func ApplyMigration(db *sqlx.DB, cfg *models.ConfigEnv, logger logger.Logger) error {
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
