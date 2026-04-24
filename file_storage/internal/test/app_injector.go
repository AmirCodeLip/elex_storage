package test

import (
	"elex_storage/file_storage/internal/domain"
	"elex_storage/file_storage/internal/infrastructure/configs"
	"elex_storage/file_storage/internal/infrastructure/database"
	"elex_storage/pkg/logger"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

type AppTest struct {
}

func setConfigs() *domain.ConfigEnv {
	cfg := configs.TestConfigEnv()
	// ToDo AmirCodelip remove logger path by changing logger test logger.
	os.Setenv("LOGGER_PATH", "D:\\Projects\\elex_storage\\bin\\test.log")
	return cfg
}

func InjectBase(t *testing.T) (*fxtest.App, *sqlx.DB, logger.Logger, *domain.ConfigEnv) {
	var db *sqlx.DB
	var lg logger.Logger
	var cfg *domain.ConfigEnv
	app := fxtest.New(
		t,
		fx.Provide(setConfigs),
		fx.Provide(logger.NewLogger),
		fx.Provide(database.NewDatabase),
		fx.Populate(&db),
		fx.Populate(&lg),
		fx.Populate(&cfg),
	)
	app.Start(t.Context())
	return app, db, lg, cfg
}

// ToDo AmirCodelip Create Mocks
// func InjectMock(t *testing.T) (*fxtest.App, logger.Logger, repositories.DirectoryMetadataRepository, repositories.FileMetadataRepository) {
// 	setConfigs()
// 	var directoryMetadataRepository repositories.DirectoryMetadataRepository
// 	var fileMetadataRepository repositories.FileMetadataRepository
// 	var lg logger.Logger
// 	app := fxtest.New(
// 		t,
// 		fx.Provide(setConfigs),
// 		fx.Provide(logger.NewLogger),
// 		fx.Provide(mock_repositories.CreateMockDirectoryMetadataRepository),
// 		fx.Provide(mock_repositories.CreateMockFileMetadataRepository),
// 		fx.Populate(&directoryMetadataRepository),
// 		fx.Populate(&fileMetadataRepository),
// 		fx.Populate(&lg),
// 	)
// 	app.Start(t.Context())
// 	return app, lg, directoryMetadataRepository, fileMetadataRepository
// }
