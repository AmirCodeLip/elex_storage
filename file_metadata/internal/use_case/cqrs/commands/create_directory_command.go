package commands

import (
	"github.com/google/uuid"
)

type CreateDirectoryCommand struct {
	Name     string     `json:"name"`
	ParentId *uuid.UUID `json:"parent_id"`
}
