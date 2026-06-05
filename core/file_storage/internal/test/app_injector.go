package test

import (
	"elex_storage/file_storage/internal/infrastructure/database"
	"elex_storage/pkg/logger"
	"elex_storage/pkg/shared_kernel"
	"elex_storage/pkg/shared_kernel/models"
	"testing"

	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

type AppTest struct {
	App    *fxtest.App
	Db     *sqlx.DB
	Logger logger.Logger
	Cfg    *models.ConfigEnv
}

func setConfigs() (*models.ConfigEnv, error) {
	cfg, err := shared_kernel.NewTestConfigEnv("file_storage")
	return cfg, err
}

func InjectBase(t *testing.T) AppTest {
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
		fx.Provide(database.NewTestDatabase),
		fx.Invoke(database.ApplyMigration),
		fx.Populate(&db),
		fx.Populate(&lg),
		fx.Populate(&cfg),
	)
	app.Start(t.Context())
	return AppTest{app, db, lg, cfg}
}

func InjectMock(t *testing.T) (*fxtest.App, logger.Logger, *models.ConfigEnv) {
	var lg logger.Logger
	var cfg *models.ConfigEnv
	app := fxtest.New(
		t,
		fx.Provide(func() *testing.T {
			return t
		}),
		fx.Provide(setConfigs),
		fx.Provide(logger.NewLoggerMock),
		fx.Populate(&lg),
		fx.Populate(&cfg),
	)
	app.Start(t.Context())
	return app, lg, cfg
}
