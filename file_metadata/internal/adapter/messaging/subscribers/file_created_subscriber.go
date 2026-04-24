package subscribers

import (
	"elex_storage/file_metadata/internal/use_case/cqrs/command_handlers"
	"elex_storage/file_metadata/internal/use_case/cqrs/commands"
	"elex_storage/pkg/logger"
	"elex_storage/pkg/shared_kernel/event_models"
	"encoding/json"
)

type FileCreatedSubscriber struct {
	fileCreatedHandler *command_handlers.FileCreatedHandler
	logger             logger.Logger
}

func NewFileCreatedSubscriber(fileCreatedHandler *command_handlers.FileCreatedHandler, logger logger.Logger) *FileCreatedSubscriber {
	return &FileCreatedSubscriber{fileCreatedHandler, logger}
}

func (s *FileCreatedSubscriber) Handle(msg []byte) error {
	var fileEntity event_models.FileCreated
	if err := json.Unmarshal([]byte(msg), &fileEntity); err == nil {
		createErr := s.fileCreatedHandler.Handle(commands.FileCreatedCommand{FileEntity: fileEntity})
		if createErr != nil {
			s.logger.Error(createErr.Error())
		}
	} else {
		s.logger.Error(err.Error())
	}
	return nil
}
