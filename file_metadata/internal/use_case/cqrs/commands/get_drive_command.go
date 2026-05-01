package commands

import "github.com/google/uuid"

type GetDriveCommand struct {
	ParentId uuid.UUID
}
