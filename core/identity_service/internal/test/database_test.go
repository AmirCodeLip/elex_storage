package test

import (
	"elex_storage/pkg/auto_migration"
	"testing"
)

func TestMigration(t *testing.T) {
	app, db, logger, cfg := InjectBase(t)

	autoMigrateManager, err := auto_migration.NewAutoMigrateManager(db, cfg.MigrationsDir)
	if err != nil {
		logger.Error(err.Error())
		t.Fatal(err.Error())
	}
	// Apply Migrations
	migrationsErr := autoMigrateManager.AutoMigrate()
	if migrationsErr != nil {
		logger.Error(migrationsErr.Error())
		t.Fatal(migrationsErr.Error())
	}
	app.Done()
}
