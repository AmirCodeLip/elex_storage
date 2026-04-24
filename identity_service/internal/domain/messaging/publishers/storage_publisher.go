package publishers

import (
	"elex_storage/pkg/shared_kernel/event_models"
)

type StoragePublisher interface {
	PublishUserRegisterd(event event_models.UserRegisterd) error
}
