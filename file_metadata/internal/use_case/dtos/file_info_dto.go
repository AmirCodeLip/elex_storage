package dtos

import (
	"time"

	"github.com/google/uuid"
)

type FileInfoDto struct {
	Id        uuid.UUID
	Name      string
	Size      int
	CreatedAt time.Time
}
