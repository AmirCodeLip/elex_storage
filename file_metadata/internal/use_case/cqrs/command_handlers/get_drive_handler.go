package command_handlers

import (
	"elex_storage/file_metadata/internal/domain/repositories"
	"elex_storage/file_metadata/internal/use_case/cqrs/commands"
	"elex_storage/file_metadata/internal/use_case/dtos"
	"elex_storage/pkg/logger"

	"github.com/jinzhu/copier"
)

type GetDriveHandler struct {
	logger                      logger.Logger
	directoryMetadataRepository repositories.DirectoryMetadataRepository
	fileMetadataRepository      repositories.FileMetadataRepository
}

func NewGetDriveHandler(logger logger.Logger,
	directoryMetadataRepository repositories.DirectoryMetadataRepository,
	fileMetadataRepository repositories.FileMetadataRepository) *GetDriveHandler {
	return &GetDriveHandler{logger, directoryMetadataRepository, fileMetadataRepository}
}

func (u *GetDriveHandler) Handle(cmd commands.GetDriveCommand) ([]dtos.StorageItemDto, error) {
	var finalResult []dtos.StorageItemDto

	// Step1: Get directories from database
	var dirsResult []dtos.StorageItemDto
	dirs, err := u.directoryMetadataRepository.GetDirectories()
	if err != nil {
		return nil, err
	}
	copier.Copy(&dirsResult, &dirs)
	for _, directoryMetadata := range dirsResult {
		directoryMetadata.Type = dtos.Directory
		finalResult = append(finalResult, directoryMetadata)
	}

	// Step2: Get files from database
	var filesResult []dtos.StorageItemDto
	files, err := u.fileMetadataRepository.GetFiles()
	if err != nil {
	}
	copier.Copy(&filesResult, &files)
	for _, fileMetadata := range filesResult {
		fileMetadata.Type = dtos.Directory
		finalResult = append(finalResult, fileMetadata)
	}

	return finalResult, nil
}
