package test

import (
	"elex_storage/identity_service/internal/infrastructure/database"
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
	return cfg
}

func InjectBase(t *testing.T) (*fxtest.App, *sqlx.DB, logger.Logger, *models.ConfigEnv) {
	var db *sqlx.DB
	var lg logger.Logger
	var cfg *models.ConfigEnv
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
