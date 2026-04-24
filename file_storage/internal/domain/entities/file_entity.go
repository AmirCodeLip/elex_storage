package entities

import (
	"elex_storage/pkg/shared_kernel/entities"
	"time"

	"github.com/google/uuid"
)

type FileEntity struct {
	Id          uuid.UUID
	Name        string
	CreatedAt   time.Time
	ContentType entities.FileContentType
}
