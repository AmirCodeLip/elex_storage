package commands

import (
	"elex_storage/pkg/shared_kernel/event_models"
)

type FileCreatedCommand struct {
	FileEntity event_models.FileCreated
}
