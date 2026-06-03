package pgx_repositories

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"elex_storage/file_metadata/internal/domain/entities"
	domain_errors "elex_storage/file_metadata/internal/domain/errors"
	"elex_storage/file_metadata/internal/domain/repositories"
	"elex_storage/pkg/logger"
	"encoding/hex"

	"github.com/jackc/pgx/v5/pgconn"
)

type DirectoryMetadataRepository struct {
	logger    logger.Logger
	db        *sqlx.DB
	driveDisk string
}

func CreateDirectoryMetadataRepository(logger logger.Logger, db *sqlx.DB) repositories.DirectoryMetadataRepository {
	return &DirectoryMetadataRepository{logger: logger, db: db}
}

func (repo *DirectoryMetadataRepository) Insert(directoryMetadataEntity *entities.DirectoryMetadataEntity) error {
	if directoryMetadataEntity.ParentId == nil || *directoryMetadataEntity.ParentId == uuid.Nil {
		// Set root for the parent id
		root, err := repo.GetRoot()
		if err != nil {
			repo.logger.Error(err.Error())
			return domain_errors.ErrInsertFileMetadata
		}
		directoryMetadataEntity.ParentId = &root.Id
	}
	directoryMetadataEntity.Id = uuid.New()
	directoryMetadataEntity.Hash = repo.Hash(directoryMetadataEntity)
	_, insertErr := repo.db.NamedExec(`INSERT INTO directories_metadata (id, hash, name, parent_id, size) 
		VALUES (:id, :hash, :name, :parent_id, :size)`, directoryMetadataEntity)
	if insertErr != nil {
		var pgErr *pgconn.PgError
		if errors.As(insertErr, &pgErr) {
			fmt.Println(pgErr.Code)
			if pgErr.Code == "23505" {
				return domain_errors.ErrDirectoryExists(directoryMetadataEntity.Name)
			}
		}
		repo.logger.Error(insertErr.Error())
		return domain_errors.ErrInsertFileMetadata
	}
	return nil
}

func (repo *DirectoryMetadataRepository) GetDirectories(parentId uuid.UUID) (*[]entities.DirectoryMetadataEntity, error) {
	var dirs []entities.DirectoryMetadataEntity
	query := "SELECT * FROM directories_metadata WHERE parent_id = $1"
	err := repo.db.Select(&dirs, query, parentId)
	return &dirs, err
}

func (repo *DirectoryMetadataRepository) GetRoot() (*entities.DirectoryMetadataEntity, error) {
	root := entities.DirectoryMetadataEntity{Name: "~"}

	err := repo.db.Get(&root, "SELECT * FROM directories_metadata WHERE name=$1 AND parent_id IS NULL", root.Name)
	if errors.Is(err, sql.ErrNoRows) {
		// Create new root.
		root.Id = uuid.New()
		root.Hash = repo.Hash(&root)
		root.ParentId = nil
		_, insertErr := repo.db.NamedExec(`INSERT INTO directories_metadata (id, hash, name, parent_id, size) 
			VALUES (:id, :hash, :name, :parent_id, :size)`, root)
		if insertErr != nil {
			return nil, insertErr
		}

	} else if err != nil {
		return nil, err
	}
	return &root, nil
}

// setDirectoryId sets the directory ID for the file, defaulting to root if nil
func (repo *DirectoryMetadataRepository) SetDirectoryId(fileMetadataEntity *entities.FileMetadataEntity) error {
	if fileMetadataEntity.DirectoryId != nil {
		return nil
	}

	root, err := repo.GetRoot()
	if err != nil {
		return fmt.Errorf("failed to get root directory: %w", err)
	}

	if root == nil {
		return errors.New("root directory not found")
	}

	fileMetadataEntity.DirectoryId = &root.Id
	return nil
}

func (repo *DirectoryMetadataRepository) Hash(directoryMetadataEntity *entities.DirectoryMetadataEntity) string {
	hasher := sha256.New()
	hasher.Write([]byte(directoryMetadataEntity.Name))
	hash := hasher.Sum(nil)
	if directoryMetadataEntity.ParentId != nil {
		hasher.Write([]byte(directoryMetadataEntity.ParentId.String()))
		hash = hasher.Sum(nil)
	}
	return hex.EncodeToString(hash)
}
