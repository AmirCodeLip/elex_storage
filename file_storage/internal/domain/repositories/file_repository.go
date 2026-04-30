package repositories

import (
	"elex_storage/file_storage/internal/domain/entities"

	"github.com/jmoiron/sqlx"
)

type FileRepository interface {
	Insert(fileEntity *entities.FileEntity, tx *sqlx.Tx) error
	GetByHash(checksum string) (*entities.FileEntity, error)
	BeginTransaction() (*sqlx.Tx, error)
}
