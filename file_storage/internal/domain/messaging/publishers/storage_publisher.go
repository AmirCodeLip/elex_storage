package publishers

import (
	"elex_storage/pkg/shared_kernel/event_models"
)

type StoragePublisher interface {
	PublishFileCreated(event event_models.FileCreated) error
}
