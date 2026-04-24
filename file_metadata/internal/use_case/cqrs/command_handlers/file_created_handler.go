package command_handlers

import (
	"elex_storage/file_metadata/internal/domain/entities"
	"elex_storage/file_metadata/internal/domain/repositories"
	"elex_storage/file_metadata/internal/use_case/cqrs"
	"elex_storage/file_metadata/internal/use_case/cqrs/commands"
	"elex_storage/pkg/logger"
	"elex_storage/pkg/shared_kernel/models"
	"errors"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type FileCreatedHandler struct {
	logger                      logger.Logger
	fileMetadataRepository      repositories.FileMetadataRepository
	directoryMetadataRepository repositories.DirectoryMetadataRepository
}

func NewFileCreatedHandler(logger logger.Logger,
	fileMetadataRepository repositories.FileMetadataRepository, directoryMetadataRepository repositories.DirectoryMetadataRepository) *FileCreatedHandler {
	return &FileCreatedHandler{logger, fileMetadataRepository, directoryMetadataRepository}
}

func (u *FileCreatedHandler) Handle(cmd commands.FileCreatedCommand) error {
	newFileId := uuid.New()
	var fileMetadataEntity entities.FileMetadataEntity
	// Map event into entity file
	copier.Copy(&fileMetadataEntity, &cmd.FileEntity)
	fileMetadataEntity.Id = newFileId
	// Set Storage id into storage_id
	fileMetadataEntity.StorageId = cmd.FileEntity.Id
	ext := filepath.Ext(fileMetadataEntity.Name)
	fileMetadataEntity.Name = strings.TrimSuffix(fileMetadataEntity.Name, ext)
	fileMetadataEntity.FileExtension = ext
	if fileMetadataEntity.DirectoryId == nil {
		root, err := u.directoryMetadataRepository.GetRoot()
		if err != nil {
			return err
		}
		fileMetadataEntity.DirectoryId = &root.Id
	}
	err := u.fileMetadataRepository.Insert(fileMetadataEntity)
	if err != nil {
		var commonErr *models.CommonError
		if errors.As(err, &commonErr) {
			return commonErr
		}
		u.logger.Info(err.Error())
		return cqrs.SaveFileMetadataErr()
	}
	return nil
}
