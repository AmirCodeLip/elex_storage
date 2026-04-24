package auto_migration

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

var migrationsErr = errors.New("Can't read and apply migreations")

type migrationInfo struct {
	created_at time.Time
	path       string
	name       string
}

type AutoMigrateManager struct {
	Path       string
	migrations []*migrationInfo
	db         *sqlx.DB
}

func NewAutoMigrateManager(db *sqlx.DB, migrationsDir string) (*AutoMigrateManager, error) {
	// Get migrations directory
	pwd, _ := os.Getwd()
	var exeDir = filepath.Dir(pwd)
	migrationDir := filepath.Join(exeDir, migrationsDir) // directory of migrations
	info, err := os.Stat(migrationDir)
	if err != nil {
		return nil, errors.Join(migrationsErr, err)
	}
	if !info.IsDir() {
		return nil, errors.Join(migrationsErr, errors.New("migrations is not directory"))
	}
	// Create migration table
	createMigrationsQuery := `
		CREATE TABLE IF NOT EXISTS migrations (
			ID VARCHAR(255) PRIMARY KEY,
			DIRTY BOOLEAN NOT NULL,
			APPLIED_AT TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`
	_, err = db.Exec(createMigrationsQuery)
	if err != nil {
		return nil, errors.Join(migrationsErr, err)
	}

	return &AutoMigrateManager{Path: migrationDir, db: db}, nil
}

func (autoMigrateManager *AutoMigrateManager) ReadFiles() error {
	var migrations []*migrationInfo
	errReadFiles := filepath.Walk(autoMigrateManager.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Join(migrationsErr, err)
		}
		// // skip directories
		if info.IsDir() {
			return nil
		}
		fileName := info.Name()
		ext := filepath.Ext(fileName)
		fileName = strings.TrimSuffix(fileName, ext)
		created_at, err := time.Parse("20060102", fileName[0:8])
		if err != nil {
			return errors.Join(migrationsErr, err)
		}

		// return nil
		fmt.Println(created_at)
		migrations = append(migrations, &migrationInfo{
			created_at: created_at,
			name:       fileName,
			path:       path,
		})
		return nil
	})
	autoMigrateManager.migrations = migrations
	return errReadFiles
}

// ToDo AmirCodelip Sort migration files

func (autoMigrateManager *AutoMigrateManager) AutoMigrate() error {
	err := autoMigrateManager.ReadFiles()
	if err != nil {
		return err
	}
	var refError error = nil
	for _, migration := range autoMigrateManager.migrations {
		sqlBytes, err := os.ReadFile(migration.path)
		if err != nil {
			refError = errors.Join(migrationsErr, err)
			break
		}
		var exists bool
		err = autoMigrateManager.db.Get(&exists, `SELECT EXISTS (SELECT 1 FROM migrations WHERE id = $1 AND dirty = false)`, migration.name)
		if err != nil {
			refError = errors.Join(migrationsErr, err)
			break
		}
		if exists {
			continue
		}
		// Apply migration file into the database
		sqlQ := string(sqlBytes)
		_, err = autoMigrateManager.db.Exec(sqlQ)
		autoMigrateManager.db.Exec(`INSERT INTO migrations(id, dirty) VALUES ($1, $2);`, migration.name, err != nil)
		if err != nil {
			refError = errors.Join(migrationsErr, err)
			break
		}
	}
	return refError
}
