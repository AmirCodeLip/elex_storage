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

func setConfigs() (*models.ConfigEnv, error) {
	test_env, err := shared_kernel.GetTestEnvPath()
	if err != nil {
		return nil, err
	}
	cfg, err := shared_kernel.TestConfigEnv(test_env)
	return cfg, err
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
