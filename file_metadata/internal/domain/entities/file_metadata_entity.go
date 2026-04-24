package entities

import (
	"time"

	"github.com/google/uuid"
)

type FileMetadataEntity struct {
	Id            uuid.UUID       `db:"id"`
	StorageId     uuid.UUID       `db:"storage_id"`
	Name          string          `db:"name"`
	FileExtension string          `db:"file_extension"`
	ContentType   FileContentType `db:"content_type"`
	Size          int             `db:"size"`
	Drive         string          `db:"drive"`
	Hash          string          `db:"hash"`
	DirectoryId   *uuid.UUID      `db:"directory_id"`
	CreatedAt     time.Time       `db:"created_at"`
	UpdateAt      time.Time       `db:"update_at"`
}
