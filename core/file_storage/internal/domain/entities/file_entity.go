package entities

import (
	"elex_storage/pkg/shared_kernel/entities"
	"time"

	"github.com/google/uuid"
)

type FileEntity struct {
	Id          uuid.UUID                `db:"id"`
	CreatedAt   time.Time                `db:"created_at"`
	ContentType entities.FileContentType `db:"content_type"`
	Checksum    string                   `db:"checksum"`
	IsMounted   string                   `db:"is_mounted"`
}
