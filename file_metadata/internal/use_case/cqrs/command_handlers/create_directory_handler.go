package command_handlers

import (
	"elex_storage/file_metadata/internal/domain/entities"
	"elex_storage/file_metadata/internal/domain/repositories"
	"elex_storage/file_metadata/internal/use_case/cqrs/commands"
	"elex_storage/pkg/logger"
)

type CreateDirectoryHandler struct {
	logger                      logger.Logger
	directoryMetadataRepository repositories.DirectoryMetadataRepository
}

func NewCreateDirectoryHandler(logger logger.Logger,
	directoryMetadataRepository repositories.DirectoryMetadataRepository) *CreateDirectoryHandler {
	return &CreateDirectoryHandler{logger, directoryMetadataRepository}
}

func (u *CreateDirectoryHandler) Handle(cmd commands.CreateDirectoryCommand) (error, entities.DirectoryMetadataEntity) {
	directory := entities.DirectoryMetadataEntity{Name: cmd.Name, ParentId: cmd.ParentId}
	err := u.directoryMetadataRepository.Insert(&directory)
	return err, directory
}
