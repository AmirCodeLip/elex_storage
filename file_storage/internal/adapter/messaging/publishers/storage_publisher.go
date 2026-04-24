package publishers

import (
	"elex_storage/file_storage/internal/domain/messaging/publishers"
	"elex_storage/pkg/message_broker"
	"elex_storage/pkg/shared_kernel/event_models"
	"encoding/json"
)

type StoragePublisher struct {
	messaging message_broker.EventMessaging
}

func CreateStoragePublisher(messaging message_broker.EventMessaging) publishers.StoragePublisher {
	return &StoragePublisher{messaging}
}

func (s *StoragePublisher) PublishFileCreated(event event_models.FileCreated) error {
	data, _ := json.Marshal(event)
	return s.messaging.SendMessage("file.created", string(data))
}
