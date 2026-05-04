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

func (u *GetDriveHandler) Handle(cmd commands.GetDriveCommand) ([]dtos.DirectoryDto, []dtos.FileDto, error) {

	// Step1: Get directories from database
	var dirsResult []dtos.DirectoryDto
	dirs, err := u.directoryMetadataRepository.GetDirectories()
	if err != nil {
		return nil, nil, err
	}
	copier.Copy(&dirsResult, &dirs)

	// Step2: Get files from database
	var filesResult []dtos.FileDto
	files, err := u.fileMetadataRepository.GetFiles()
	if err != nil {
	}
	copier.Copy(&filesResult, &files)

	return dirsResult, filesResult, nil
}
