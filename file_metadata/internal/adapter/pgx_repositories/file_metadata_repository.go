package pgx_repositories

import (
	"crypto/sha256"
	"elex_storage/file_metadata/internal/domain/entities"
	domain_errors "elex_storage/file_metadata/internal/domain/errors"
	"elex_storage/file_metadata/internal/domain/repositories"
	"elex_storage/pkg/logger"
	"encoding/hex"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

type FileMetadataRepository struct {
	logger    logger.Logger
	db        *sqlx.DB
	driveDisk string
}

func CreateFileMetadataRepository(logger logger.Logger, db *sqlx.DB) repositories.FileMetadataRepository {
	return &FileMetadataRepository{logger: logger, db: db}
}

func (repo *FileMetadataRepository) GetFiles() (*[]entities.FileMetadataEntity, error) {
	var files []entities.FileMetadataEntity
	err := repo.db.Select(&files, "SELECT * FROM files_metadata")
	return &files, err
}

func (repo *FileMetadataRepository) Insert(fileMetadataEntity entities.FileMetadataEntity) error {
	if fileMetadataEntity.DirectoryId == nil {
		return errors.New("DirectoryId can't be null")
	}
	fileMetadataEntity.Hash = repo.Hash(fileMetadataEntity)
	_, insertErr := repo.db.NamedExec(`INSERT INTO files_metadata(id, storage_id, name, file_extension, content_type, size, drive, hash, directory_id) 
		VALUES (:id, :storage_id, :name, :file_extension, :content_type, :size, :drive, :hash, :directory_id)`,
		fileMetadataEntity,
	)
	if insertErr != nil {
		var pgErr *pgconn.PgError
		if errors.As(insertErr, &pgErr) {
			if pgErr.Code == "23505" {
				return domain_errors.ErrFileExists
			}
		}
		repo.logger.Error(insertErr.Error())
		return domain_errors.ErrInsertFileMetadata
	}
	return nil
}

func (repo *FileMetadataRepository) Hash(fileMetadataEntity entities.FileMetadataEntity) string {
	hasher := sha256.New()
	hasher.Write([]byte(fileMetadataEntity.DirectoryId.String()))
	hasher.Write([]byte(fileMetadataEntity.Name))
	hasher.Write([]byte(fileMetadataEntity.FileExtension))
	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash)
}
