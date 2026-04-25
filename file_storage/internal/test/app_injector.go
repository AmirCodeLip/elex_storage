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
