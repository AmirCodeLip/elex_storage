package entities

import (
	"time"

	"github.com/google/uuid"
)

type DirectoryMetadataEntity struct {
	Id        uuid.UUID  `db:"id"`
	Name      string     `db:"name"`
	Hash      string     `db:"hash"`
	Size      int        `db:"size"`
	ParentId  *uuid.UUID `db:"parent_id"`
	CreatedAt time.Time  `db:"created_at"`
}
