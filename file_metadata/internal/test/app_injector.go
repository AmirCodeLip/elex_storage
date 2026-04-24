package test

import (
	"elex_storage/file_metadata/internal/domain/repositories"
	"elex_storage/file_metadata/internal/infrastructure/database"
	mock_repositories "elex_storage/file_metadata/internal/test/repositories"
	"elex_storage/pkg/logger"
	"elex_storage/pkg/shared_kernel"
	"elex_storage/pkg/shared_kernel/models"
	"testing"

	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

type AppTest struct {
}

func setConfigs() *models.ConfigEnv {
	cfg := shared_kernel.TestConfigEnv()
	// os.Setenv("LOGGER_PATH", "D:\\Projects\\elex_storage\\bin\\test.log")
	return cfg
}

func InjectBase(t *testing.T) (*fxtest.App, *sqlx.DB, logger.Logger) {
	var db *sqlx.DB
	var lg logger.Logger
	app := fxtest.New(
		t,
		fx.Provide(func() *testing.T {
			return t
		}),
		fx.Provide(setConfigs),
		fx.Provide(logger.NewLoggerMock),
		fx.Provide(database.NewDatabase),
		fx.Populate(&db),
		fx.Populate(&lg),
	)
	app.Start(t.Context())
	return app, db, lg
}

func InjectMock(t *testing.T) (*fxtest.App, logger.Logger, repositories.DirectoryMetadataRepository, repositories.FileMetadataRepository) {
	setConfigs()
	var directoryMetadataRepository repositories.DirectoryMetadataRepository
	var fileMetadataRepository repositories.FileMetadataRepository
	var lg logger.Logger
	app := fxtest.New(
		t,
		fx.Provide(func() *testing.T {
			return t
		}),
		fx.Provide(setConfigs),
		fx.Provide(logger.NewLoggerMock),
		fx.Provide(mock_repositories.CreateMockDirectoryMetadataRepository),
		fx.Provide(mock_repositories.CreateMockFileMetadataRepository),
		fx.Populate(&directoryMetadataRepository),
		fx.Populate(&fileMetadataRepository),
		fx.Populate(&lg),
	)
	app.Start(t.Context())
	return app, lg, directoryMetadataRepository, fileMetadataRepository
}
