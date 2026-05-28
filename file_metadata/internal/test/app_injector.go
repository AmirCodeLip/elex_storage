package test

import (
	"elex_storage/file_metadata/internal/domain/repositories"
	"elex_storage/file_metadata/internal/infrastructure/database"
	mock_repositories "elex_storage/file_metadata/internal/test/repositories"
	"elex_storage/pkg/logger"
	"elex_storage/pkg/shared_kernel"
	"elex_storage/pkg/shared_kernel/models"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

type AppTest struct {
}

func setConfigs() (*models.ConfigEnv, error) {
	pwd, _ := os.Getwd()
	var exeDir = filepath.Dir(pwd)
	paths := strings.Split(exeDir, "\\")
	for {
		/// remove until get root of project
		last := paths[len(paths)-1]
		if last != "elex_storage" {
			paths = paths[:len(paths)-1]
		} else {
			break
		}
	}
	paths = append(paths, ".env")
	envAddr := strings.Join(paths, "\\")
	return shared_kernel.TestConfigEnv(&envAddr)
}

func InjectBase(t *testing.T) (*fxtest.App, *sqlx.DB, logger.Logger) {
	var db *sqlx.DB
	var lg logger.Logger
	app := fxtest.New(
		t,
		fx.Provide(func() *testing.T {
			return t
		}),
		fx.Provide(logger.NewLoggerMock),
		fx.Provide(setConfigs),
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
