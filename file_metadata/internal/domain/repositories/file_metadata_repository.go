package repositories

import (
	"elex_storage/file_metadata/internal/domain/entities"

	"github.com/jmoiron/sqlx"
)

type FileMetadataRepository interface {
	BeginTransaction() (*sqlx.Tx, error)
	RollbackTransaction(tx *sqlx.Tx) error
	CommitTransaction(tx *sqlx.Tx) error
	Insert(fileMetadataEntity entities.FileMetadataEntity, tx *sqlx.Tx) error
	GetFiles() (*[]entities.FileMetadataEntity, error)
}
