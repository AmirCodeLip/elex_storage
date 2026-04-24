package test

import (
	"elex_storage/pkg/auto_migration"
	"testing"
)

func TestMigration(t *testing.T) {
	app, db, _, cfg := InjectBase(t)

	autoMigrateManager, err := auto_migration.NewAutoMigrateManager(db, cfg.MigrationsDir)
	if err != nil {
		t.Fatal(err.Error())
	}
	// Apply Migrations
	migrationsErr := autoMigrateManager.AutoMigrate()
	if migrationsErr != nil {
		t.Fatal(err.Error())
	}
	app.Done()
}
