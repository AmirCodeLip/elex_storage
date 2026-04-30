package event_models

import (
	"elex_storage/pkg/shared_kernel/entities"

	"github.com/google/uuid"
)

type FileCreated struct {
	Id          uuid.UUID                `json:"id"`
	Name        string                   `json:"name"`
	ContentType entities.FileContentType `json:"content-type"`
	Size        int                      `json:"size"`
	Drive       string                   `json:"drive"`
	Checksum    string                   `json:"checksum"`
	Key         []byte                   `json:"key"`
}
