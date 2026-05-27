package command_handlers

import (
	"elex_storage/file_metadata/internal/domain/repositories"
	"elex_storage/file_metadata/internal/use_case/cqrs/commands"
	"elex_storage/file_metadata/internal/use_case/dtos"
	"elex_storage/pkg/logger"

	"github.com/google/uuid"
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
	// Step: Set parentId and if it was null, get root dir.
	var parentId uuid.UUID = uuid.Nil
	if cmd.ParentId == nil {
		root, err := u.directoryMetadataRepository.GetRoot()
		if err != nil {
			return nil, nil, err
		}
		parentId = root.Id
	} else {
		/// ToDo AmirCodelip check parrent id exist.
		parentId = *cmd.ParentId
	}

	// Step1: Get directories from database
	var dirsResult []dtos.DirectoryDto
	dirs, err := u.directoryMetadataRepository.GetDirectories(parentId)
	if err != nil {
		return nil, nil, err
	}
	copier.Copy(&dirsResult, &dirs)

	// Step2: Get files from database
	var filesResult []dtos.FileDto
	files, err := u.fileMetadataRepository.GetFiles(parentId)
	if err != nil {
	}
	copier.Copy(&filesResult, &files)

	return dirsResult, filesResult, nil
}
